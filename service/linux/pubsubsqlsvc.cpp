#include <cstdlib>
#include <unistd.h>
#include <sys/wait.h>
#include <syslog.h>
#include <stdio.h>
#include <string.h>
#include <string>
#include <iostream>

int install(const char* serviceFile, const std::string& options);
int uninstall();
void rundaemon(char* path, char** argv);
void* logthread(void *);

// wrap "must succeed" system calls into os class to avoid error handling spaghetti
class os {
public:
	static void pipe(int fildes[2], const char* desc);
	static int fork(const char* desc);
	static void close(int fd, const char* desc);
	static void dup2(int oldfd, int newfd, const char* desc);
	static void execvp(const char *file, char *const argv[], const char* desc);
	static void pthread_create(pthread_t *thread, const pthread_attr_t *attr, void *(*start_routine) (void *), void *arg, const char* desc);
	static void waitpid(pid_t pid, int *status, int options, const char* desc);
	static FILE* fdopen(int fd, const char *mode, const char* desc);
	static void setsid(const char* desc);
	static void chdir(const char* path, const char* desc);
private:
	static void logfatal(const char* err, const char* desc);
};

int logpipe[2];
int inpipe[2]; 
const char* logprefix = "pubsubsql";

int main(int argc, char** argv) {
	// validate command line input
	std::string usage = " valid commands [install, uninstall]";
	if (argc < 2) {
		std::cerr << "invalid command: " <<  usage << std::endl;
		return EXIT_FAILURE;
	}
	// execute
	std::string command(argv[1]);
	if (command == "install") { 
		std::string options;
		return install(argv[0], options);
	} else if (command == "uninstall") { 
		return uninstall();
	} else if (command == "svc") { 
		if (argc < 3) {
			std::cerr << "expected executable file path" << std::endl;
			return EXIT_FAILURE;
		}
		rundaemon(argv[2], argv + 2);
	}
	// invalid command
	std::cerr << "invalid command: " << command << usage << std::endl;
	return EXIT_FAILURE;
}

int install(const char* serviceFile, const std::string& options) {
	std::cout << "install" << std::endl;
	return EXIT_SUCCESS;
}

int uninstall() {
	std::cout << "uninstall" << std::endl;
	return EXIT_SUCCESS;
}

void rundaemon(char* path, char** argv) {
	// first become a daemon
	if (os::fork("become background process") > 0) _exit(EXIT_SUCCESS); 	
	os::setsid("become a leader of new session");	
	if (os::fork("ensure we are not session leader") > 0) _exit(EXIT_SUCCESS); 	
	//os::chdir("/", "change to root directory");
	os::close(STDIN_FILENO, "closing stdin");	
	os::close(STDOUT_FILENO, "closing stdout");	
	os::close(STDERR_FILENO, "closing stderr");	
	// try to close all open file descriptors
	int maxfd = sysconf(_SC_OPEN_MAX);	
	if (-1 == maxfd) maxfd = 8192; // limit is indeterminate, so take a guess
	for (int fd = 0; fd < maxfd; fd++) close(fd);
	// we are daemon, start pubsubsql child process
	os::pipe(logpipe, "create pipe to redirect stderr to syslog");
	os::pipe(inpipe, "create pipe to redirect to stdin");
	int status = 0;
	pid_t childPid = os::fork("forking child process");
	switch (childPid) {
	case 0:  // child process context 
		os::close(logpipe[0], "close read end of the logpipe");
		os::close(inpipe[1], "close write end of inpipe");
		// associate pipe with stdin
		if (inpipe[0] != STDIN_FILENO) {
			os::dup2(inpipe[0], STDIN_FILENO, "duplicate read pipe->STDIN_FILENO");
			os::close(inpipe[0], "close duplicated read pipe->STDIN_FILENO");
		}
		// associate pipe with stderr
		if (logpipe[1] != STDERR_FILENO) {
			os::dup2(logpipe[1], STDERR_FILENO, "duplicate write pipe->STDERR_FILENO");
			os::close(logpipe[1], "close duplicated write pipe->STDERR_FILENO");
		}
		os::execvp(path, argv, "starting pubsubsql");
		break;
	default: // parent process context
		os::close(logpipe[1], "close write end of the logpipe");
		os::close(inpipe[0], "close read end of inpipe");
		// start thread to redirect err from pubsubsql to syslog			
		pthread_t thread;
		os::pthread_create(&thread, NULL, logthread, NULL, "start logger thread");
		os::waitpid(childPid, &status, 0, "waiting for child pubsubsql");
		pthread_join(thread, NULL);
		_exit(status);
	}		
	_exit(EXIT_SUCCESS);
}

void* logthread(void *) {
	FILE* f = os::fdopen(logpipe[0], "r", "set stream in logthread");	
	openlog(logprefix, LOG_PERROR, LOG_USER);
	const int BUFFER_SIZE = 4096;	
	char buffer[1 + BUFFER_SIZE] = {0};
	for (;;) {
		const char* line = fgets(buffer, BUFFER_SIZE, f); 
		if (NULL == line) {
			// if we fail to read it indicates that child process is done
			break;
		}
		// redirect log message to syslog
		if (strncmp(line, "info", 4) == 0) {
			syslog(LOG_INFO, "%s", line);
		} else if (strncmp(line, "error", 5) == 0) {
			syslog(LOG_ERR, "%s", line);
		} else if (strncmp(line, "debug", 5) == 0) {
			syslog(LOG_DEBUG, "%s", line);
		} else {
			syslog(LOG_WARNING, "%s", line);
		}
	}
	closelog();
	return NULL;
}

// os
void os::pipe(int fildes[2], const char* desc) {
	if (-1 == ::pipe(fildes)) logfatal("pipe() failed", desc);
}

int os::fork(const char* desc) {
	int ret = ::fork();
	if (-1 == ret) logfatal("fork() failed", desc);
	return ret;
}

void os::close(int fd, const char* desc) {
	if (-1 == ::close(fd)) logfatal("close() failed", desc);		
}

void os::dup2(int oldfd, int newfd, const char* desc) {
	if (-1 == ::dup2(oldfd, newfd)) logfatal("dup2() failed", desc);	
}

void os::execvp(const char *file, char *const argv[], const char* desc) {
	::execvp(file, argv);
	logfatal("execvp() failed", desc);	
}

void os::pthread_create(pthread_t *thread, const pthread_attr_t *attr, void *(*start_routine) (void *), void *arg, const char* desc) {
	if (0 != ::pthread_create(thread, attr, start_routine, arg)) logfatal("pthread_create() failed", desc);
}

void os::waitpid(pid_t pid, int *status, int options, const char* desc) {
	if (-1 == ::waitpid(pid, status, options)) logfatal("waitpid() failed", desc);
}

FILE* os::fdopen(int fd, const char *mode, const char* desc) {
	FILE* ret = ::fdopen(fd, mode);
	if (NULL == ret) logfatal("fdopen() failed", desc);
	return ret;
}

void os::setsid(const char* desc) {
	if (-1 == ::setsid()) logfatal("setsid() failed", desc);
}

void os::chdir(const char* path, const char* desc) {
	if (-1 == ::chdir(path)) logfatal("chdir() failed", desc);
}

// since we are not running in a terminal session, log to syslog and exit the process
void os::logfatal(const char* err, const char* desc) {
	openlog(logprefix, LOG_PID | LOG_PERROR, LOG_USER);
	syslog(LOG_EMERG, "%s %s", err, desc);
	closelog();
	_exit(EXIT_FAILURE);
}


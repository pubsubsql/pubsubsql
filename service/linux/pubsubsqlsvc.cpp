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
int runAsDaemon(char* path, char** argv);
void* logthread(void *);

int logpipe[2];
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
		return runAsDaemon(argv[2], argv + 2);
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

int runAsDaemon(char* path, char** argv) {
	// create pipe to redirect stderr to syslog
	pipe(logpipe);
	int status = 0;
	pid_t childPid = fork();
	//
	switch (childPid) {
	case -1:
		break;
	case 0: // start pubsubsql
		// close read end of the logpipe
		close(logpipe[0]);
		// associate pipe with stderr
		if (logpipe[1] != STDERR_FILENO) {
			dup2(logpipe[1], STDERR_FILENO);
			close(logpipe[1]);
		}
		execvp(path, argv);
		break;
	default:
		// close write end of the logpipe
		close(logpipe[1]);
		// redirect err from pubsubsql to syslog			
		pthread_t tx;
		pthread_create(&tx, NULL, logthread, NULL);
		if (waitpid(childPid, &status, 0) == -1) {
			break;
		} 
		pthread_join(tx, NULL);
		break;
	}		

	return EXIT_SUCCESS;
}

void* logthread(void *) {
	openlog(logprefix, LOG_PERROR, LOG_USER);
	FILE* f = fdopen(logpipe[0], "r");	
	if (NULL == f) {
		return NULL;
	} 
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


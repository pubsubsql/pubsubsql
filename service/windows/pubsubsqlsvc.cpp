/* Copyright (C) 2013 CompleteDB LLC.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with PubSubSQL.  If not, see <http://www.gnu.org/licenses/>.
 */

#include <iostream>
#include <thread>
#include "process.h"
#include "eventlog.h"

int install(const char* serviceFile, const std::string& options);
int uninstall();
int runAsService();

const char* SERVICE_NAME = "PubSubSQL";

int main(int argc, char* argv[])
{
	// validate command line input
	std::string usage = "valid commands [install, uninstall]";
	if (argc < 2) {
		std::cerr << "no command found: " << usage << std::endl;
		return EXIT_FAILURE;
	}
	// assemble options
	std::string options;
	for (int i = 2; i < argc; i++) {
		char* arg = argv[i];
		options.append(" ");
		options.append(arg);
	}
	std::cout << options << std::endl;
	// execute
	std::string command(argv[1]);
	if (command == "install") return install(argv[0], options);
	else if (command == "uninstall") return uninstall();
	else if (command == "svc") return runAsService();
	// invalid command
	std::cerr << "invalid command: " << usage << std::endl;
	return EXIT_FAILURE;
	/*
	pipe testPipe;
	if (!testPipe.ok()) {
		std::cerr << "pipe failed" << std::endl;
	}
	//
	std::thread t([&]() {
		std::cout<< testPipe.readLine();
	});
	testPipe.writeLine("pipe test");
	t.join();
	// test filepath
	std::cout << eventlog::getPath() << std::endl;
	*/
	//eventlog::install("pubsubsqlsvc.exe", "pubsubsql");
	/*
	eventlog log("pubsubsql");
	log.logdebug("debug");
	log.loginfo("info");
	log.logwarn("warn");
	log.logerror("error");
	*/
	// test process redirection

	/*
	process pubsubsql;
	if (pubsubsql.start("C:\\Users\\Oleg\\Go\\src\\pubsubsql\\pubsubsql.exe")) {
		int i = 0;
		std::cin >> i;
		pubsubsql.stop();
		pubsubsql.wait(3000);
	}
	*/
	return 0;
}

int install(const char* serviceFile, const std::string& options) {
	int ret = EXIT_SUCCESS;
	std::cout << "installing " << SERVICE_NAME << " service" << std::endl;
	SC_HANDLE manager = OpenSCManager(NULL, SERVICES_ACTIVE_DATABASE, SC_MANAGER_CREATE_SERVICE);
	if (NULL == manager) {
		std::cout << "failed to connect to service control manager error:" << GetLastError() << std::endl;
		return EXIT_FAILURE;
	}
	std::string servicePath;
	servicePath.append(serviceFile);
	servicePath.append(" svc ");
	servicePath.append(options);
	SC_HANDLE service = CreateService(manager, SERVICE_NAME, SERVICE_NAME, SERVICE_START | SERVICE_STOP | DELETE,
		SERVICE_WIN32_OWN_PROCESS, SERVICE_DEMAND_START, SERVICE_ERROR_NORMAL, servicePath.c_str(),
		NULL, NULL, NULL, NULL, NULL);
	if (NULL == service) {
		std::cout << "failed to install service error:" << GetLastError() << std::endl;	
		ret = EXIT_FAILURE;
	}
	CloseServiceHandle(manager);
	CloseServiceHandle(service);
	//
	//
	if (ret != EXIT_FAILURE) {
		eventlog::install("pubsubsqlsvc.exe", "pubsubsql");
		std::cout << "service " << SERVICE_NAME << " was installed " << std::endl;
	} else {
		std::cout << "MAKE SURE YOU ARE RUNNING WITH REQUIRED SECURITY PRIVILEGES AND SERVICE DOES NOT ALREADY EXIST!" << std::endl;
	}
	return ret;
}

int uninstall() {
	int ret = EXIT_SUCCESS;
	std::cout << "uninstalling " << SERVICE_NAME << " service" << std::endl;
	SC_HANDLE manager = OpenSCManager(NULL, SERVICES_ACTIVE_DATABASE, SC_MANAGER_CREATE_SERVICE);
	if (NULL == manager) {
		std::cout << "failed to connect to service control manager error:" << GetLastError() << std::endl;
		return EXIT_FAILURE;
	}
	SC_HANDLE service = OpenService(manager, SERVICE_NAME, SC_MANAGER_ALL_ACCESS);
	if (NULL == service) {
		std::cout << "failed to open service error:" << GetLastError() << std::endl;	
		ret = EXIT_FAILURE;
	}
	if (!DeleteService(service)) {
		std::cout << "failed to uninstall service error:" << GetLastError() << std::endl;	
		ret = EXIT_FAILURE;
	}
	CloseServiceHandle(manager);
	CloseServiceHandle(service);
	if (ret != EXIT_FAILURE) {
		std::cout << "service " << SERVICE_NAME << " was uninstalled " << std::endl;
	} else {
		std::cout << "MAKE SURE YOU ARE RUNNING WITH REQUIRED SECURITY PRIVILEGES AND SERVICE EXISTS!" << std::endl;
	}
	return ret;
}

int runAsService() {
	return EXIT_SUCCESS;
}

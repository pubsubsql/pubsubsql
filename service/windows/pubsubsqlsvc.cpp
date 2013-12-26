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

int main(int argc, char* argv[])
{
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
	eventlog::install("pubsubsqllog.dll", "pubsubsql");
	eventlog log("pubsubsql");
	log.logdebug("debug");
	log.loginfo("info");
	log.logwarn("warn");
	log.logerror("error");
	// test process redirection
	process pubsubsql;
	if (pubsubsql.start("C:\\Users\\Oleg\\Go\\src\\pubsubsql\\pubsubsql.exe")) {
		int i = 0;
		std::cin >> i;
		pubsubsql.stop();
		pubsubsql.wait(3000);
	}
	return 0;
}


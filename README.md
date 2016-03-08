[![Build Status](https://travis-ci.org/pubsubsql/pubsubsql.svg?branch=master)](https://travis-ci.org/pubsubsql/pubsubsql) [![GitHub version](https://badge.fury.io/gh/pubsubsql%2Fpubsubsql.svg)](https://badge.fury.io/gh/pubsubsql%2Fpubsubsql) 

pubsubsql
=========

Homepage: [http://pubsubsql.com/](http://pubsubsql.com/)


An open-source in-memory database offering:
  - SQL-like syntax
  - PUB-SUB functionality
  - optional streaming to MySQL

Developers usually have to choose between a fast in-memory store and a conventional database to give their users a responsive real-time feeling and at the same time allow backend service run conventional analytics in a more analytic-friendly environment. We've been in that same boat and ended up coding complex contraptions with caches, long-term storage and notification systems which more often than not spiral out of control and become a maintenance nightmare. So we thought "Why not have a SQL-like interface with a PUB-SUB interface and just stream data to a normal database?". Result of that work is PubSubSQL! The project has been dormant for some time but we've regrouped and decided to see where it will take us.

PubSubSQL is fully written in Go and supports a multitude of clients - Go, Node.js, Python, Java...

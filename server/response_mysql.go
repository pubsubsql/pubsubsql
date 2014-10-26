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

package server

//=====================================================================================================================
// cmdMysqlConnectResponse
//---------------------------------------------------------------------------------------------------------------------
type cmdMysqlConnectResponse struct {
	requestIdResponse
	address string
	error string
}

func newCmdMysqlConnectResponse(req *mysqlConnectRequest) *cmdMysqlConnectResponse {
	return &cmdMysqlConnectResponse {
		address: req.address,
		error: "",
	}
}

func (this *cmdMysqlConnectResponse) toNetworkReadyJSON() ([]byte, bool) {
	builder := networkReadyJSONBuilder()
	builder.beginObject()
	ok(builder)
	builder.valueSeparator()
	action(builder, "mysqlConnect")
	builder.valueSeparator()
	builder.nameValue("address", this.address)
	if "" != this.error {
		builder.valueSeparator()
		builder.nameValue("error", this.error)
	}
	builder.endObject()
	return builder.getNetworkBytes(this.requestId), false
}
//=====================================================================================================================
// cmdMysqlDisconnectResponse
//---------------------------------------------------------------------------------------------------------------------
type cmdMysqlDisconnectResponse struct {
	requestIdResponse
	error string
}

func newCmdMysqlDisconnectResponse(req *mysqlDisconnectRequest) *cmdMysqlDisconnectResponse {
	return &cmdMysqlDisconnectResponse {
		error: "",
	}
}

func (this *cmdMysqlDisconnectResponse) toNetworkReadyJSON() ([]byte, bool) {
	builder := networkReadyJSONBuilder()
	builder.beginObject()
	ok(builder)
	builder.valueSeparator()
	action(builder, "mysqlDisconnect")
	if "" != this.error {
		builder.valueSeparator()
		builder.nameValue("error", this.error)
	}
	builder.endObject()
	return builder.getNetworkBytes(this.requestId), false
}
//=====================================================================================================================
// cmdMysqlStatusResponse
//---------------------------------------------------------------------------------------------------------------------
type cmdMysqlStatusResponse struct {
	requestIdResponse
	online int
	error string
}

func newCmdMysqlStatusResponse(req *mysqlStatusRequest) *cmdMysqlStatusResponse {
	return &cmdMysqlStatusResponse {
		online: 0,
		error: "",
	}
}

func (this *cmdMysqlStatusResponse) toNetworkReadyJSON() ([]byte, bool) {
	builder := networkReadyJSONBuilder()
	builder.beginObject()
	ok(builder)
	builder.valueSeparator()
	action(builder, "mysqlStatus")
	builder.valueSeparator()
	builder.nameIntValue("online", this.online)
	if "" != this.error {
		builder.valueSeparator()
		builder.nameValue("error", this.error)
	}
	builder.endObject()
	return builder.getNetworkBytes(this.requestId), false
}

func (this *cmdMysqlStatusResponse) setOnline(online bool) {
	if online {
		this.online = 1
	} else {
		this.online = 0
	}
}

func (this *cmdMysqlStatusResponse) isOnline() (bool) {
	return (this.online != 0)
}

func (this *cmdMysqlStatusResponse) isOffline() (bool) {
	return ! this.isOnline()
}

//=====================================================================================================================
// cmdMysqlSubscribeResponse
//---------------------------------------------------------------------------------------------------------------------
type cmdMysqlSubscribeResponse struct {
	requestIdResponse
	error string
}

func newCmdMysqlSubscribeResponse(req *mysqlSubscribeRequest) *cmdMysqlSubscribeResponse {
	return &cmdMysqlSubscribeResponse {
		error: "",
	}
}

func (this *cmdMysqlSubscribeResponse) toNetworkReadyJSON() ([]byte, bool) {
	builder := networkReadyJSONBuilder()
	builder.beginObject()
	ok(builder)
	builder.valueSeparator()
	action(builder, "mysqlSubscribe")
	if "" != this.error {
		builder.valueSeparator()
		builder.nameValue("error", this.error)
	}
	builder.endObject()
	return builder.getNetworkBytes(this.requestId), false
}
//=====================================================================================================================
// cmdMysqlUnsubscribeResponse
//---------------------------------------------------------------------------------------------------------------------
type cmdMysqlUnsubscribeResponse struct {
	requestIdResponse
	error string
}

func newCmdMysqlUnsubscribeResponse(req *mysqlUnsubscribeRequest) *cmdMysqlUnsubscribeResponse {
	return &cmdMysqlUnsubscribeResponse {
		error: "",
	}
}

func (this *cmdMysqlUnsubscribeResponse) toNetworkReadyJSON() ([]byte, bool) {
	builder := networkReadyJSONBuilder()
	builder.beginObject()
	ok(builder)
	builder.valueSeparator()
	action(builder, "mysqlUnsubscribe")
	if "" != this.error {
		builder.valueSeparator()
		builder.nameValue("error", this.error)
	}
	builder.endObject()
	return builder.getNetworkBytes(this.requestId), false
}
//=====================================================================================================================
// cmdMysqlTablesResponse
//---------------------------------------------------------------------------------------------------------------------
type cmdMysqlTablesResponse struct {
	requestIdResponse
	tables []string
	error string
}

func newCmdMysqlTablesResponse(req *mysqlTablesRequest) *cmdMysqlTablesResponse {
	return &cmdMysqlTablesResponse {
		tables : make([]string, 0),
		error: "",
	}
}

func (this *cmdMysqlTablesResponse) toNetworkReadyJSON() ([]byte, bool) {
	builder := networkReadyJSONBuilder()
	builder.beginObject()
	ok(builder)
	builder.valueSeparator()
	action(builder, "mysqlTables")
	builder.valueSeparator()
	builder.string("tables")
	builder.nameSeparator()
	builder.beginArray()
	for i, tableName := range this.tables {
		// another tableName
		if i != 0 {
			builder.valueSeparator()
		}
		builder.string(tableName)
	}
	builder.endArray()
	if "" != this.error {
		builder.valueSeparator()
		builder.nameValue("error", this.error)
	}
	builder.endObject()
	return builder.getNetworkBytes(this.requestId), false
}
//=====================================================================================================================

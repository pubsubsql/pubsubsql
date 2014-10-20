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
	connectionAddress string
}

func newCmdMysqlConnectResponse(req *mysqlConnectRequest) *cmdMysqlConnectResponse {
	return &cmdMysqlConnectResponse {
		connectionAddress: req.connectionAddress,
	}
}

func (this *cmdMysqlConnectResponse) toNetworkReadyJSON() ([]byte, bool) {
	builder := networkReadyJSONBuilder()
	builder.beginObject()
	ok(builder)
	builder.valueSeparator()
	action(builder, "mysqlConnect")
	builder.valueSeparator()
	builder.nameValue("connectionAddress", this.connectionAddress)
	builder.endObject()
	return builder.getNetworkBytes(this.requestId), false
}
//=====================================================================================================================
// cmdMysqlDisconnectResponse
//---------------------------------------------------------------------------------------------------------------------
type cmdMysqlDisconnectResponse struct {
	requestIdResponse
}

func newCmdMysqlDisconnectResponse(req *mysqlDisconnectRequest) *cmdMysqlDisconnectResponse {
	return &cmdMysqlDisconnectResponse {
		// void
	}
}

func (this *cmdMysqlDisconnectResponse) toNetworkReadyJSON() ([]byte, bool) {
	builder := networkReadyJSONBuilder()
	builder.beginObject()
	ok(builder)
	builder.valueSeparator()
	action(builder, "mysqlDisconnect")
	builder.endObject()
	return builder.getNetworkBytes(this.requestId), false
}
//=====================================================================================================================
// cmdMysqlStatusResponse
//---------------------------------------------------------------------------------------------------------------------
type cmdMysqlStatusResponse struct {
	requestIdResponse
	connectionOnline int
}

func newCmdMysqlStatusResponse(req *mysqlStatusRequest) *cmdMysqlStatusResponse {
	return &cmdMysqlStatusResponse {
		connectionOnline: 0,
	}
}

func (this *cmdMysqlStatusResponse) toNetworkReadyJSON() ([]byte, bool) {
	builder := networkReadyJSONBuilder()
	builder.beginObject()
	ok(builder)
	builder.valueSeparator()
	action(builder, "mysqlStatus")
	builder.valueSeparator()
	builder.nameIntValue("connectionOnline", this.connectionOnline)
	builder.endObject()
	return builder.getNetworkBytes(this.requestId), false
}

func (this *cmdMysqlStatusResponse) setOnline() {
	this.connectionOnline = 1
}

func (this *cmdMysqlStatusResponse) setOffline() {
	this.connectionOnline = 0
}

func (this *cmdMysqlStatusResponse) isOnline() (bool) {
	return (this.connectionOnline > 0)
}

func (this *cmdMysqlStatusResponse) isOffline() (bool) {
	return ! this.isOnline()
}

//=====================================================================================================================
// cmdMysqlSubscribeResponse
//---------------------------------------------------------------------------------------------------------------------
type cmdMysqlSubscribeResponse struct {
	requestIdResponse
}

func newCmdMysqlSubscribeResponse(req *mysqlSubscribeRequest) *cmdMysqlSubscribeResponse {
	return &cmdMysqlSubscribeResponse {
		// void
	}
}

func (this *cmdMysqlSubscribeResponse) toNetworkReadyJSON() ([]byte, bool) {
	builder := networkReadyJSONBuilder()
	builder.beginObject()
	ok(builder)
	builder.valueSeparator()
	action(builder, "mysqlSubscribe")
	builder.endObject()
	return builder.getNetworkBytes(this.requestId), false
}
//=====================================================================================================================
// cmdMysqlUnsubscribeResponse
//---------------------------------------------------------------------------------------------------------------------
type cmdMysqlUnsubscribeResponse struct {
	requestIdResponse
}

func newCmdMysqlUnsubscribeResponse(req *mysqlUnsubscribeRequest) *cmdMysqlUnsubscribeResponse {
	return &cmdMysqlUnsubscribeResponse {
		// void
	}
}

func (this *cmdMysqlUnsubscribeResponse) toNetworkReadyJSON() ([]byte, bool) {
	builder := networkReadyJSONBuilder()
	builder.beginObject()
	ok(builder)
	builder.valueSeparator()
	action(builder, "mysqlUnsubscribe")
	builder.endObject()
	return builder.getNetworkBytes(this.requestId), false
}
//=====================================================================================================================

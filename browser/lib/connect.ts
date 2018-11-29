import {IConnectionFactory} from "./types";

export class ConnectionFactory implements IConnectionFactory{
    constructor(
        private readonly __config: RTCConfiguration,
        private readonly __table: { [key: string]: (RTCPeerConnection) => any } = {}){

    }
    create(channel: string = "data"): RTCPeerConnection {
        const conn = new RTCPeerConnection(this.__config);
        const init = this.__table[channel];
        if  (init && !init(conn)){
            conn.close();
            return null;
        }
        return conn;
    }
}

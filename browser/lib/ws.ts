import {IDriver, IMessage} from "./types";

export class WebSocketDriver implements IDriver{
    private __ws: WebSocket;
    private readonly __table: Map<string,any> = new Map<string, any>();
    constructor(
        private readonly __id: string,
        private readonly __url: string,
    ){
        const ws = new WebSocket(`${this.__url}/${this.__id}`);
        ws.onmessage = (e) => this.__message(JSON.parse(e.data) as IMessage,ws);
        this.__ws = ws;
    }
    private async __message(data: IMessage,ws?: WebSocket) {
        const val = this.__table.get(data.command);
        if(!!val){
            val(data.payload,this,data.source);
        }
    }
    public on(command:string, callback: any) :IDriver{
        this.__table.set(command,callback);
        return this;
    }
    public emit(command: string,payload: any): IDriver {
        this.__ws.send(JSON.stringify({
            source: this.__id,
            payload,
            command,
        }));
        return this;
    }
    public emitTo(id: string,command: string,payload: any): IDriver {
        this.__ws.send(JSON.stringify({
            source: this.__id,
            dest:id,
            payload,
            command,
        }));
        return this;
    }
}

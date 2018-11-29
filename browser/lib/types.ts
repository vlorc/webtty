
export interface IMessage{
    command: string;
    dest?: string;
    source?: string;
    payload: any;
}

export interface IDescription{
    session: string;
    channel?: string;
    type?: RTCSdpType;
    sdp: string;
}

export interface ISession extends RTCPeerConnection{
    id(): string;
    state(): Promise<number>;
}

export interface ISessionFactory{
    id(): string;
    create(id: string,conn: RTCPeerConnection): ISession;
    attach(id: string,conn: RTCPeerConnection): ISession;
    query(id: string): ISession;
    remove(id: string): void;
}

export interface IConnectionFactory{
    create(channel: string): RTCPeerConnection;
}

export interface IDriver{
    on(command:string, any) :IDriver;
    emit(command: string,payload: any): IDriver;
    emitTo(id: string,command: string,payload: any): IDriver;
}

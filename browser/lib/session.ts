import { ISessionFactory, ISession } from "./types";

const tables = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXTZabcdefghiklmnopqrstuvwxyz".split("");

export class SessionFactory implements ISessionFactory {
    private readonly __session: Record<string, ISession> = {};
    constructor(private readonly __id: string) {

    }

    public id(): string {
        return this.__id;
    }

    public attach(id: string,conn: RTCPeerConnection): ISession{
        (conn as any).state = function (): Promise<number>{
            return new Promise(function (resolve,reject) {
                this.__resolve = resolve;
            }.bind(this));
        };
        (conn as any).id = function () {
            return id;
        };
        const __close = conn.close.bind(conn);
        (conn as any).close = () => {
            this.remove(id);
            __close();
        };
        this.__session[id] = conn as ISession;
        return conn as ISession;
    }

    public create(id: string,conn: RTCPeerConnection): ISession {
        return this.attach(this.__random(id),conn);
    }

    public query(id: string): ISession {
        return this.__session[id];
    }

    public remove(id: string): void {
        delete (this.__session[id]);
    }

    public __random(id: string, size: number = 8): string {
        if (!length) {
            length = Math.floor(Math.random() * tables.length);
        }
        let str = "";
        for (let i = 0; i < size; i++) {
            str += tables[Math.floor(Math.random() * tables.length)];
        }
        return [this.__id, id, str].join(".");
    }
}


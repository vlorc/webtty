import {IMessage, IDescription, ISessionFactory, ISession, IConnectionFactory, IDriver} from "./types";

export default class Peer{
    public message: any;
    constructor(
        private readonly __connect: IConnectionFactory,
        private readonly __session: ISessionFactory,
        private readonly __driver: IDriver){

        this.__driver.on("OFFER", (v,_,s) => this.__onOffer(v as IDescription,s));
        this.__driver.on("ANSWER", (v,_,s) => this.__onAnswer(v as IDescription));
        this.__driver.on("CLOSE", (v,_,s) => this.__onClose(v as IDescription));
    }

    public async connect(id: string,channel: string = "data",option: any = null,...init: Array<(RTCPeerConnection) => any>): Promise<ISession> {
        let conn = this.__connect.create(channel);
        for(let v of init){
            if(!v(conn)){
                conn.close();
                return null;
            }
        }
        
        let session = this.__session.create(id,conn);
        let offer = await session.createOffer(option);
        await session.setLocalDescription(offer);
        return new Promise<ISession>((resolve,reject) => {
            session.onicecandidate = e => {
                if (e.candidate === null) {
                    this.__driver.emitTo(id,"OFFER",{
                        channel,
                        session: session.id(),
                        type: session.localDescription.type,
                        sdp: session.localDescription.sdp,
                    });
                    resolve(session);
                }
            };
        });
    }

    private __create(des: IDescription): ISession{
        let session = this.__session.query(des.session);
        if(!session){
            session = this.__session.attach(des.session,this.__connect.create(des.channel))
        }
        return session;
    }

    private async __onOffer(des: IDescription,source: string){
        const session = this.__create(des);
        await session.setRemoteDescription(new RTCSessionDescription({
            type: "offer",
            sdp: des.sdp,
        }));

        const answer = await session.createAnswer();
        await session.setLocalDescription(answer);

        return new Promise<ISession>((resolve,reject) => {
            session.onicecandidate = e => {
                if (e.candidate === null) {
                    this.__driver.emitTo(source,"ANSWER",{
                        session: des.session,
                        channel: des.channel,
                        type: session.localDescription.type,
                        sdp: session.localDescription.sdp,
                    });
                    resolve();
                }
            };
        });
    }

    private async __onClose(des: IDescription){
        const session = this.__session.query(des.session);
        if(!session){
            session.close();
        }
    }

    private async __onAnswer(des: IDescription){
        const session = this.__session.query(des.session);
        await session.setRemoteDescription(new RTCSessionDescription({
            type: "answer",
            sdp: des.sdp,
        }));

        if((session as any).__resolve){
            (session as any).__resolve(1);
        }
    }
}

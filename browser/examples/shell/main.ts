
import {SessionFactory, ConnectionFactory, Peer, WebSocketDriver} from "peer";
import { Terminal } from 'xterm';
import { fit } from 'xterm/lib/addons/fit/fit';

const peer = open("qnmb");

function log(msg: string) {
    console.log(msg);
}

function open(id: string) {
    const conn = new Peer(
        new ConnectionFactory({
                iceServers: [
                    {
                        urls: "turn:www.xxxxxxxx.cn:3478",
                        username: "guest",
                        credential: "guest",
                    }
                ]
            }),
        new SessionFactory(id),
        new WebSocketDriver(id,`wss://api.xxxxxxx.cn/peer`),
    );
    return conn;
}

function shell(conn: RTCPeerConnection) {
    conn.onsignalingstatechange = e => log(`state: ${conn.signalingState}`);
    conn.oniceconnectionstatechange = e => log(`change: ${conn.iceConnectionState}`);
    const dc = conn.createDataChannel("data");
    const cc = conn.createDataChannel("control");
    conn.ondatachannel = function (e: RTCDataChannelEvent) {
        const channel = e.channel;
        log(`channel event: '${channel.label}'`);
        channel.onclose = () => log(`channel close: '${channel.label}'`);
        channel.onopen = () => log(`channel open:  '${channel.label}'`);
        channel.onmessage = (e) => {
            log(`channel message: '${channel.label}' payload '${e.data}'`);
        }
    };
    dc.onclose = () => console.log('sendChannel has closed')
    dc.onopen = () => {
        const term = new Terminal({
            cols: 200,
            rows: 60,
        });
        term.open(document.getElementById("shell"));
        fit(term);
        term.on('title', function (title) {
            document.title = title;
        });
        term.on('resize', function ({ cols, rows }) {
            cc.send(JSON.stringify({ type: "resize", "data": [cols, rows] }));
        });
        term.on('data', function (data) {
            dc.send(data);
        });
        dc.onmessage = e => term.write(e.data);
    };
    cc.onmessage = e => log(`Message from DataChannel '${cc.label}' payload '${e.data}'`);
    return true;
}

async function connect(id: string) {
    const session = await peer.connect(id, "shell", null, shell);
    const state = await session.state();
    log("RTCPeerConnection open finish");
}

(window as any).startSession = () => {
    let remote = (document.getElementById('remoteId') as HTMLInputElement).value;
    if (remote === '') {
        return alert('Session Description must not be empty')
    }
    connect(remote);
};

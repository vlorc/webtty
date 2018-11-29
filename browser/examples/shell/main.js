"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : new P(function (resolve) { resolve(result.value); }).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
var __generator = (this && this.__generator) || function (thisArg, body) {
    var _ = { label: 0, sent: function() { if (t[0] & 1) throw t[1]; return t[1]; }, trys: [], ops: [] }, f, y, t, g;
    return g = { next: verb(0), "throw": verb(1), "return": verb(2) }, typeof Symbol === "function" && (g[Symbol.iterator] = function() { return this; }), g;
    function verb(n) { return function (v) { return step([n, v]); }; }
    function step(op) {
        if (f) throw new TypeError("Generator is already executing.");
        while (_) try {
            if (f = 1, y && (t = op[0] & 2 ? y["return"] : op[0] ? y["throw"] || ((t = y["return"]) && t.call(y), 0) : y.next) && !(t = t.call(y, op[1])).done) return t;
            if (y = 0, t) op = [op[0] & 2, t.value];
            switch (op[0]) {
                case 0: case 1: t = op; break;
                case 4: _.label++; return { value: op[1], done: false };
                case 5: _.label++; y = op[1]; op = [0]; continue;
                case 7: op = _.ops.pop(); _.trys.pop(); continue;
                default:
                    if (!(t = _.trys, t = t.length > 0 && t[t.length - 1]) && (op[0] === 6 || op[0] === 2)) { _ = 0; continue; }
                    if (op[0] === 3 && (!t || (op[1] > t[0] && op[1] < t[3]))) { _.label = op[1]; break; }
                    if (op[0] === 6 && _.label < t[1]) { _.label = t[1]; t = op; break; }
                    if (t && _.label < t[2]) { _.label = t[2]; _.ops.push(op); break; }
                    if (t[2]) _.ops.pop();
                    _.trys.pop(); continue;
            }
            op = body.call(thisArg, _);
        } catch (e) { op = [6, e]; y = 0; } finally { f = t = 0; }
        if (op[0] & 5) throw op[1]; return { value: op[0] ? op[1] : void 0, done: true };
    }
};
Object.defineProperty(exports, "__esModule", { value: true });
var peer_1 = require("peer");
var xterm_1 = require("xterm");
var fit_1 = require("xterm/lib/addons/fit/fit");
var peer = open("qnmb");
function log(msg) {
    // document.getElementById('logs').innerHTML += msg + '<br>'
}
function open(id) {
    var conn = new peer_1.Peer(new peer_1.ConnectionFactory({
        iceServers: [
            {
                urls: "turn:www.atatakai.cn:3478",
                username: "guest",
                credential: "guest",
            }
        ]
    }), new peer_1.SessionFactory(id), new peer_1.WebSocketDriver(id, "wss://api.atatakai.cn/peer"));
    return conn;
}
function shell(conn) {
    conn.onsignalingstatechange = function (e) { return log("state: " + conn.signalingState); };
    conn.oniceconnectionstatechange = function (e) { return log("change: " + conn.iceConnectionState); };
    var dc = conn.createDataChannel("data");
    var cc = conn.createDataChannel("control");
    conn.ondatachannel = function (e) {
        var channel = e.channel;
        log("channel event: '" + channel.label + "'");
        channel.onclose = function () { return log("channel close: '" + channel.label + "'"); };
        channel.onopen = function () { return log("channel open:  '" + channel.label + "'"); };
        channel.onmessage = function (e) {
            log("channel message: '" + channel.label + "' payload '" + e.data + "'");
        };
    };
    dc.onclose = function () { return console.log('sendChannel has closed'); };
    dc.onopen = function () {
        var term = new xterm_1.Terminal({
            cols: 200,
            rows: 60,
        });
        term.open(document.getElementById("shell"));
        fit_1.fit(term);
        term.on('title', function (title) {
            document.title = title;
        });
        term.on('resize', function (_a) {
            var cols = _a.cols, rows = _a.rows;
            cc.send(JSON.stringify({ type: "resize", "data": [cols, rows] }));
        });
        term.on('data', function (data) {
            dc.send(data);
        });
        dc.onmessage = function (e) { return term.write(e.data); };
    };
    cc.onmessage = function (e) { return log("Message from DataChannel '" + cc.label + "' payload '" + e.data + "'"); };
    return true;
}
function connect(id) {
    return __awaiter(this, void 0, void 0, function () {
        var session, state;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0: return [4 /*yield*/, peer.connect(id, "shell", null, shell)];
                case 1:
                    session = _a.sent();
                    return [4 /*yield*/, session.state()];
                case 2:
                    state = _a.sent();
                    log("RTCPeerConnection open finish");
                    return [2 /*return*/];
            }
        });
    });
}
window.startSession = function () {
    var remote = document.getElementById('remoteId').value;
    if (remote === '') {
        return alert('Session Description must not be empty');
    }
    connect(remote);
};

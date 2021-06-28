data = {
    ws: null, // Our websocket
    newMsg: '', // Holds new messages to be sent to the server
    chatContent: '', // A running list of chat messages displayed on the screen
    id: "",
    name: "",
    joined: false // True if email and username have been filled in
}

window.onload = function () {
    var self = data

    data.chatContent = document.getElementById("chat-space").innerHTML

    self.ws = new WebSocket("ws://" + window.location.host + "/ws")
    self.ws.addEventListener("message", function (e) {
        var msg = JSON.parse(e.data)

        //自动滑动到最下方
        var element = document.getElementById("chat-space");
        element.innerHTML += "<li>" + msg.name + "  :=>  " + msg.message + "</li>"
        element.scrollTop = element.scrollHeight
    })

    var id = ""
    let req = new XMLHttpRequest()
    req.open("GET", "/uuid");
    req.onload = function () {
        if (req.status != 200) {
            alert("获取Uid失败")
        } else {
            data.id = req.responseText
        }
    }
    req.send()
}

function sendMessage() {
    var send = document.getElementById("message-send-port")

    if (data.id == "") {

        alert("no uuid")
    }
    var message = send.getElementsByClassName("msg")[0]
    var name = send.getElementsByClassName("name")[0]
    //注册
    if (!data.joined) {
        name.setAttribute("hidden", true);
        send.getElementsByClassName("name-box")[0].setAttribute("hidden", true)
        data.name = (name.value)

        var greet = document.body.getElementsByClassName("greet")[0]
        greet.textContent = "你好" + name.value + "! 欢迎加入聊天室"
        greet.removeAttribute("hidden");

        data.joined = true
    }

    msg = message.value

    if (msg != "") {
        data.ws.send(
            JSON.stringify({
                uuid: data.id,
                name: data.name,
                message: msg
            })
        );
    }
    message.value = ""

}
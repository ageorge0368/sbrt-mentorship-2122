const username_Input = document.getElementById('username'), 
      chat_Input_Typing = document.getElementById('input_typing'),
      chat_Input = document.getElementById('chat_input'),
      chat_Log = document.getElementById('actual_chat_log');

function loginHandeler() {
    fetch(`login?${new URLSearchParams({
        username: username_Input.value
    })}`)
    .then(response => response.text())
    .then(message => {
        if(isNaN(parseInt(message))) {
            alert(message);
        } else {
            const session = new URLSearchParams({username: username_Input.value, sessionID: message}),
                  ws = new WebSocket(`wss://${location.hostname}:${location.port}/chat?${session}`);

            ws.onopen = () => alert(`Connected to backend!`);
            ws.onclose = () => alert('Disconnected from backend! :(');
            ws.onmessage = ({data}) => {
                const message = JSON.parse(data),
                      element = document.createElement(`li`);
                element.innerText = `${message.sender}: ${message.content}`;
                chat_Log.appendChild(element);
            };
            chat_Input_Typing.onsubmit = () => {
                ws.send(JSON.stringify({
                    content: chat_Input.value    
                }));
                return false;
            };
        }
    });
    return false;
}
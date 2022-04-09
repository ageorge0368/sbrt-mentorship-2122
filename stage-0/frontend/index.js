const arr = [1,2,3,34,45,65,67,7,7,88,0];
while(arr.length) {
    console.log(arr[arr.length-1]);
    arr.pop();
}



const userNameInput = document.getElementById('username'),
      chatInputForm = document.getElementById('chat-input-form'),
      chatInput = document.getElementById('chat-input'),
      chatLog = document.getElementById('chat-log'),

function loginHandle () {
    fetch(`login?${new URLSearchParams({

    })}`)
    .then(response => {
        if(response.status)
    })

}
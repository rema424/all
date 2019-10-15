'use strict';

document.addEventListener('DOMContentLoaded', () => {
  // Element
  const sendBtn = document.querySelector('#button-addon2');
  const messageInput = document.querySelector('#msgInput');
  const messages = document.querySelector('.messages');
  const messageTmpl = document.querySelector('.message-tmpl').firstElementChild;

  // Data
  let socket;

  // Event
  sendBtn.addEventListener('click', e => {
    const msg = messageInput.value;

    if (!msg) {
      e.preventDefault();
      return;
    }

    if (!socket) {
      alert('エラー：WebSocket接続が行われていません。');
      e.preventDefault();
      return;
    }

    // socket.send(msg)
    socket.send(JSON.stringify({ Body: msg }));
    messageInput.value = null;
    e.preventDefault();
    return;
  });

  // Init
  if (!window.WebSocket) {
    alert('エラー：WebSocketに対応していないブラウザです。');
  } else {
    const path = window.location.pathname.split('/');
    const roomId = path[2];
    const userId = path[3];
    socket = new WebSocket(`ws://localhost:8080/rooms/${roomId}/${userId}/socket`);

    socket.onclose = () => {
      alert('接続が終了しました。');
    };

    socket.onmessage = e => {
      console.log(e.data);
      const msg = JSON.parse(e.data);
      const frg = document.createDocumentFragment();
      frg.appendChild(messageTmpl.cloneNode(true));
      frg.querySelector('.sender small').textContent = `${msg.User || 'unknown'} : ${msg.When || ''}`;
      frg.querySelector('.message').textContent = msg.Body || 'no content';
      messages.insertBefore(frg, messages.firstElementChild);
    };
  }
});

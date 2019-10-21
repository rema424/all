document.addEventListener('DOMContentLoaded', () => {
  const box = document.querySelector('.card-body');
  if (!box) {
    return;
  }
  box.scrollTo(0, box.scrollHeight);
  box.style.visibility = 'visible';
});

// 'use strict';

// document.addEventListener('DOMContentLoaded', () => {
//   // Element
//   const chatBox = document.querySelector('#chatbox');
//   const messageInput = document.querySelector('#msgInput');
//   const messagesBox = document.querySelector('#messages');

//   // Data
//   let socket;

//   // Event
//   chatBox.addEventListener('submit', e => {
//     const msg = messageInput.value;

//     if (!msg) {
//       e.preventDefault();
//       return;
//     }

//     if (!socket) {
//       console.log('エラー：WebSocket接続が行われていません。');
//       e.preventDefault();
//       return;
//     }

//     socket.send(msg);
//     // socket.send(JSON.stringify({ Message: msg }));
//     messageInput.value = null;
//     e.preventDefault();
//     return;
//   });

//   // Func
//   const addMsgElm = msg => {
//     const li = document.createElement('li');
//     li.textContent = msg;

//     // const strong = document.createElement('strong');
//     // strong.textContent = msg.Name + ': ';
//     // li.appendChild(strong);

//     // const span = document.createElement('span');
//     // span.textContent = msg.Message;
//     // li.appendChild(span);

//     messagesBox.appendChild(li);
//   };

//   // Init
//   if (!window.WebSocket) {
//     console.log('エラー：WebSocketに対応していないブラウザです。');
//   } else {
//     const protocol = window.location.protocol == 'http:' ? 'ws:' : 'wss:';
//     const host = window.location.host;
//     const pathname = window.location.pathname;
//     socket = new WebSocket(`${protocol}//${host}/ws${pathname}`);

//     socket.onopen = () => {
//       console.log('接続を開始しました。');
//     };

//     socket.error = () => {
//       console.log('エラーが発生しました。');
//     };

//     socket.onclose = () => {
//       console.log('接続が終了しました。');
//     };

//     socket.onmessage = e => {
//       // const msg = JSON.parse(e.data);
//       const msg = e.data;
//       addMsgElm(msg);
//     };
//   }

//   let timerId;
//   let topMsgClass = 'b';

//   const handler = e => {
//     if (timerId) {
//       return;
//     }
//     timerId = setTimeout(() => {
//       timerId = 0;
//     }, 500);
//     const target = document.querySelector(`.${topMsgClass}`).getBoundingClientRect().bottom;
//     if (0 < target) {
//       console.log('scroll in', topMsgClass);
//       if (topMsgClass == 'a') {
//         window.removeEventListener('scroll', handler);
//       } else {
//         topMsgClass = 'a';
//       }
//     }
//   };

//   window.addEventListener('scroll', handler, { passive: true });
// });

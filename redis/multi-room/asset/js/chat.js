'use strict';

document.addEventListener('DOMContentLoaded', () => {
  // Element
  const box = document.querySelector('.card-body');
  if (!box) {
    return;
  }
  const msgBox = document.querySelector('.direct-chat-messages');
  const msgTmpl = document.querySelector('#msg-tmpl');
  const msgIn = document.querySelector('#msg-in');
  const sendBtn = document.querySelector('#msg-send');

  // data
  let timerId;
  let topMsgId;
  let socket;

  // Func
  const createMsgElm = msg => {
    const elm = document.createDocumentFragment();
    elm.appendChild(document.importNode(msgTmpl.content, true));
    elm.querySelector('.direct-chat-name').textContent = msg.name;
    elm.querySelector('.direct-chat-timestamp').textContent = msg.createdAt;
    elm.querySelector('.direct-chat-text').textContent = msg.body;
    elm.querySelector('.direct-chat-img').setAttribute('src', `//api.adorable.io/avatars/128/${msg.name}.png`);
    return elm;
  };

  const createMsgsElm = msgs => {
    const elm = document.createDocumentFragment();
    if (!Array.isArray(msgs)) {
      return elm;
    }
    msgs.forEach(msg => {
      elm.appendChild(createMsgElm(msg));
    });
    return elm;
  };

  const addMsgBottom = msg => {
    const elm = createMsgElm(msg);
    msgBox.appendChild(elm);
  };

  const addMsgsTop = msgs => {
    const elm = createMsgsElm(msgs);
    msgBox.firstElementChild.before(elm);
  };

  const fetchOldMsgs = topMsgId => {
    return fetch();
  };

  const scrollHandler = e => {
    console.log('a');
    if (timerId) return;

    timerId = setTimeout(() => {
      timerId = 0;
    }, 500);

    console.log(timerId);
    const elm = document.querySelector(`#msg-${topMsgId}`);
    if (!elm) return;

    const targetY = elm.getBoundingClientRect().bottom;
    if (targetY < 0) return;

    console.log('scroll in', topMsgId);
    // fetchOldMsgs(topMsgId)
    //   .then(msgs => {
    //     if (!msgs) return;

    //     topMsgId = msgs[0].id;
    //     addMsgsTop(msgs);
    //   })
    //   .catch(err => console.error(err))
    // .finally(() => console.log('f'));
  };

  // Event
  box.addEventListener('scroll', scrollHandler, { passive: true });

  sendBtn.addEventListener('click', e => {
    e.preventDefault();

    const msg = msgIn.value;

    if (!msg) {
      e.preventDefault();
      return;
    }

    if (!socket) {
      console.log('エラー：WebSocket接続が行われていません。');
      e.preventDefault();
      return;
    }

    socket.send(msg);
    // socket.send(JSON.stringify({ Message: msg }));
    msgIn.value = null;
  });

  // Init
  const m = { name: 'grace', createdAt: '1234-12-31 11:11', body: 'こんにちは' };
  const ms = Array.from({ length: 5 }, () => m);
  console.log(m);
  console.log(ms);
  addMsgBottom(m);
  addMsgsTop(ms);
  box.scrollTo(0, box.scrollHeight);
  box.style.visibility = 'visible';
  if (!window.WebSocket) {
    console.log('エラー：WebSocketに対応していないブラウザです。');
  } else {
    const protocol = window.location.protocol == 'http:' ? 'ws:' : 'wss:';
    const host = window.location.host;
    const pathname = window.location.pathname;
    socket = new WebSocket(`${protocol}//${host}/ws${pathname}`);

    socket.onopen = () => {
      console.log('接続を開始しました。');
    };

    socket.error = () => {
      console.log('エラーが発生しました。');
    };

    socket.onclose = () => {
      console.log('接続が終了しました。');
    };

    socket.onmessage = e => {
      console.log(e.data);
      // const msg = JSON.parse(e.data);
      const msg = e.data;
      // addMsgElm(msg);
    };
  }
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

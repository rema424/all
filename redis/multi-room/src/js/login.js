'use strict';

document.addEventListener('DOMContentLoaded', () => {
  const form = document.forms.namedItem('login-form');
  const submitBtn = document.querySelector('#submit');

  submitBtn.addEventListener('click', e => {
    e.preventDefault();

    const fd = new FormData(form);
    const email = fd.get('email');
    const pw = fd.get('password');

    if (email == '' || pw == '') {
      console.log('値が入力されていません。');
      return;
    }

    fd.append('remember', document.querySelector('#remember').checked);

    let status;
    fetch('/api/login', { method: 'POST', body: fd })
      .then(res => {
        status = res.status;
        return res.json();
      })
      .then(body => {
        console.log(JSON.stringify(body));
        if (status == 200) {
          // window.location.replace('/');
          console.log('ログインに成功しました。');
        }
      });
  });
});

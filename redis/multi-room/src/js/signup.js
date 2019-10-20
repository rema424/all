'use strict';

document.addEventListener('DOMContentLoaded', () => {
  const form = document.forms.namedItem('signup-form');
  const submitBtn = document.querySelector('#submit');

  submitBtn.addEventListener('click', e => {
    e.preventDefault();

    const fd = new FormData(form);
    const email = fd.get('email');
    const pw = fd.get('password');
    const pwc = fd.get('password-confirm');

    if (email == '' || pw == '' || pwc == '') {
      console.log('値が入力されていません。');
      return;
    }

    if (pw != pwc) {
      console.log('パスワードが一致しません。');
      return;
    }

    let status;
    fetch('/api/signup', { method: 'POST' })
      .then(res => {
        status = res.status;
        return res.json();
      })
      .then(body => {
        console.log(JSON.stringify(body));
        if (status == 200) {
          // window.location.replace('/');
          console.log('サインアップに成功しました。');
        }
      });
  });
});

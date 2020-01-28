# Guestbook

Guestbook is an example application showing basic usage of Google App
Engine. Users can read & write text messages and optionally log-in with
their Google account. Messages are stored in App Engine (NoSQL)
High Replication Datastore (HRD) and retrieved using a strongly consistent
(ancestor) query.

## Products

- [App Engine][1]

## Language

- [Python][2]

## APIs

- [NDB Datastore API][3]
- [Users API][4]

## Dependencies

- [webapp2][5]
- [jinja2][6]
- [Twitter Bootstrap][7]

[1]: https://developers.google.com/appengine
[2]: https://python.org
[3]: https://developers.google.com/appengine/docs/python/ndb/
[4]: https://developers.google.com/appengine/docs/python/users/
[5]: http://webapp-improved.appspot.com/
[6]: http://jinja.pocoo.org/docs/
[7]: http://twitter.github.com/bootstrap/

## E2E Test for this sample app

A Makefile is provided to deploy and run the e2e test.

To run:

     export GAE_PROJECT=your-project-id
     make

To manually run, install the requirements

    pip install -r e2e/requirements-dev.txt

Set the environment variable to point to your deployed app:

    export GUESTBOOK_URL="http://guestbook-test-dot-useful-temple-118922.appspot.com/"

Finally, run the test

    python e2e/test_e2e.py

## setup

```sh
pip install sqlalchemy pymysql
mysql.server start
mysql -uroot -e 'create database if not exists py_sqlalchemy_sandbox;'
mysql -uroot -e 'create user if not exists tester@localhost identified by "Passw0rd!";'
mysql -uroot -e 'grant all privileges on py_sqlalchemy_sandbox.* to tester@localhost;'
```

## set up

```sh
brew install python@2 python
brew upgrade python@2 python
python2 --version
pip2 --version
python3 --version
pip3 --version
pip2 install --upgrade virtualenv
pip3 install --upgrade virtualenv
cd python/appengine-guestbook-python/
ls
virtualenv --python python2 env
source /Users/rm/dev/all/python/appengine-guestbook-python/env/bin/activate
/Users/rm/dev/all/python/appengine-guestbook-python/env/bin/python2.7 -m pip install -U "pylint<2.0.0"
source env/bin/activate
ls
ls
python --version
```

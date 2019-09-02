# -*- coding: utf-8 -*-

# Define your item pipelines here
#
# Don't forget to add your pipeline to the ITEM_PIPELINES setting
# See: https://docs.scrapy.org/en/latest/topics/item-pipeline.html

import MySQLdb


class CartunePipeline(object):
    def __init__(self):
        self.conn = MySQLdb.connect(
            user='scraper',
            passwd='Passw0rd!',
            db='scrapy',
            host='127.0.0.1',
            charset='utf8',
            use_unicode=True
        )
        self.cursor = self.conn.cursor()

    def process_item(self, item, spider):
        print("process_item")
        try:
            self.cursor.execute(
                """INSERT INTO notes (url, car, user) VALUES (%s, %s, %s)""",
                (item['url'], item['car'], item['user'])
            )
            self.conn.commit()
        except MySQLdb.Error as e:
            print("Error %d: %s" % (e.args[0], e.args[1]))

        return item

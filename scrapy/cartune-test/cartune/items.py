# -*- coding: utf-8 -*-

# Define here the models for your scraped items
#
# See documentation in:
# https://docs.scrapy.org/en/latest/topics/items.html

import scrapy


class CartuneItem(scrapy.Item):
    # define the fields for your item here like:
    # name = scrapy.Field()
    # pass
    car = scrapy.Field()
    url = scrapy.Field()


class PickupDetailItem(scrapy.Item):
    car = scrapy.Field()
    url = scrapy.Field()
    user = scrapy.Field()
    tags = scrapy.Field()


class YahooJapanDetailItem(scrapy.Item):
    headline = scrapy.Field()
    url = scrapy.Field()
    title = scrapy.Field()


class PopularCarArticleItem(scrapy.Item):
    car = scrapy.Field()
    article_title = scrapy.Field()
    article_url = scrapy.Field()

# -*- coding: utf-8 -*-
import scrapy


class CarsSpider(scrapy.Spider):
    name = 'cars'
    allowed_domains = ['cartune.me']
    start_urls = ['http://cartune.me/']

    def parse(self, response):
        pass

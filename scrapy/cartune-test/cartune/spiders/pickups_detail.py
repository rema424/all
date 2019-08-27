# -*- coding: utf-8 -*-
import scrapy
from cartune.items import PickupDetailItem


class PickupsDetailSpider(scrapy.Spider):
    name = 'pickups_detail'
    allowed_domains = ['cartune.me']
    start_urls = ['http://cartune.me/']

    url_format = 'https://cartune.me{}'

    def parse(self, response):
        for pickup in response.css('div.pickupList div.pickupList-item'):
            item = PickupDetailItem()
            item['car'] = pickup.css('p a::text').extract_first()

            url = self.url_format.format(
                pickup.css('a::attr(href)').extract_first())
            requext = scrapy.Request(url=url, callback=self.parse_detail)
            requext.meta['item'] = item
            yield requext

    def parse_detail(self, response):
        item = response.meta['item']
        item['url'] = response.url
        item['user'] = response.css(
            'p.noteDetail-user-name a::text').extract_first()
        item['tags'] = response.css(
            'div.noteDetail-header ul li span::text').extract()
        yield item

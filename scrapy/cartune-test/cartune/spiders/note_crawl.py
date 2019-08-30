# -*- coding: utf-8 -*-
import scrapy
from scrapy.linkextractors import LinkExtractor
from scrapy.spiders import CrawlSpider, Rule
from cartune.items import NoteItem


class NoteCrawlSpider(CrawlSpider):
    name = 'note_crawl'
    allowed_domains = ['cartune.me']
    start_urls = ['http://cartune.me/notes/iDtrZYPwtk']

    rules = (
        Rule(LinkExtractor(restrict_css='div.noteList--card div.noteList-item',
                           allow=r'notes/'), callback='parse_item'),
        Rule(LinkExtractor(restrict_css='div.noteDetail-header ul.tagCloud--simple li')),
        Rule(LinkExtractor(restrict_css='div.noteList--card div.noteList-item')),
    )

    def parse_item(self, response):
        item = NoteItem()
        item['url'] = response.url
        item['car'] = response.css(
            'p.noteDetail-user-car a::text').extract_first()
        item['user'] = response.css(
            'p.noteDetail-user-name a::text').extract_first()
        item['tags'] = response.css(
            'div.noteDetail-header ul span::text').extract()
        return item

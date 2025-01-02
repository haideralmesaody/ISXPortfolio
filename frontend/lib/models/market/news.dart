class NewsItem {
  final String title;
  final String link;
  final String date;
  final String ticker;
  final List<String> attachments;

  NewsItem({
    required this.title,
    required this.link,
    required this.date,
    this.ticker = '',
    this.attachments = const [],
  });

  factory NewsItem.fromJson(Map<String, dynamic> json) {
    return NewsItem(
      title: json['title'] as String,
      link: json['link'] as String,
      date: json['date'] as String,
      ticker: json['ticker'] as String? ?? '',
      attachments: (json['attachments'] as List<dynamic>?)
          ?.map((e) => e['url'] as String)
          .toList() ?? [],
    );
  }
} 
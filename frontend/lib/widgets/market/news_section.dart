import 'package:flutter/material.dart';
import '../../models/market/news.dart';
import '../../services/market/news_service.dart';
import '../../services/logger_service.dart';
import 'package:url_launcher/url_launcher.dart';

class NewsSection extends StatelessWidget {
  final NewsService _newsService = NewsService();

  NewsSection({super.key});

  @override
  Widget build(BuildContext context) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const Text(
              'Market News',
              style: TextStyle(
                fontSize: 20,
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 16),
            FutureBuilder<List<NewsItem>>(
              future: _newsService.getMarketNews(),
              builder: (context, snapshot) {
                if (snapshot.hasError) {
                  return Text('Error: ${snapshot.error}');
                }

                if (!snapshot.hasData) {
                  return const CircularProgressIndicator();
                }

                final news = snapshot.data!;
                return ListView.builder(
                  shrinkWrap: true,
                  itemCount: news.length,
                  itemBuilder: (context, index) {
                    final item = news[index];
                    return ListTile(
                      title: Text(item.title),
                      subtitle: Text('${item.date} ${item.ticker}'),
                      trailing: item.attachments.isNotEmpty
                          ? const Icon(Icons.attachment)
                          : null,
                      onTap: () async {
                        final url = Uri.parse(item.link);
                        if (await canLaunchUrl(url)) {
                          await launchUrl(url);
                        }
                      },
                    );
                  },
                );
              },
            ),
          ],
        ),
      ),
    );
  }
} 
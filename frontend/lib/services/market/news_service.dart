import 'package:http/http.dart' as http;
import 'dart:convert';
import '../../models/market/news.dart';
import '../logger_service.dart';

class NewsService {
  final String baseUrl = 'http://localhost:8000/api';

  Future<List<NewsItem>> getMarketNews() async {
    try {
      final response = await http.get(Uri.parse('$baseUrl/market/news'));
      
      if (response.statusCode == 200) {
        final List<dynamic> newsJson = json.decode(response.body);
        return newsJson.map((json) => NewsItem.fromJson(json)).toList();
      } else {
        throw Exception('Failed to load market news');
      }
    } catch (e) {
      LoggerService.error('Error fetching market news: $e');
      rethrow;
    }
  }
} 
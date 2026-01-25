import 'dart:async';
import 'dart:convert';
import 'package:flutter/foundation.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:http/http.dart' as http;
import '../../../core/api/api_client.dart';

final chatRepositoryProvider =
    Provider((ref) => ChatRepository(ref.read(apiClientProvider)));

class ChatMessage {
  final String role; // 'user' or 'assistant'
  final String content;
  final DateTime timestamp;

  ChatMessage({
    required this.role,
    required this.content,
    DateTime? timestamp,
  }) : timestamp = timestamp ?? DateTime.now();

  Map<String, dynamic> toJson() => {
        'role': role,
        'content': content,
      };
}

class ChatRepository {
  final ApiClient _client;

  ChatRepository(this._client);

  /// GET /ai/chat/suggestions - Get suggested questions
  Future<List<String>> getSuggestions() async {
    try {
      final response = await _client.get('/ai/chat/suggestions');
      final data = response.data;
      if (data is Map && data.containsKey('suggestions')) {
        return List<String>.from(data['suggestions']);
      }
      return [];
    } catch (e) {
      debugPrint('Error fetching suggestions: $e');
      return [
        "How do I identify Fall Armyworm in maize?",
        "What are the best organic pesticides for tomatoes?",
        "When is the best time to plant cassava in Uganda?",
      ];
    }
  }

  /// POST /ai/chat - Send message (non-streaming)
  Future<Map<String, dynamic>> sendMessage(
    String message, {
    String? context,
    List<ChatMessage>? history,
  }) async {
    try {
      final response = await _client.post('/ai/chat', data: {
        'message': message,
        if (context != null) 'context': context,
        if (history != null) 'history': history.map((m) => m.toJson()).toList(),
      });
      return response.data as Map<String, dynamic>;
    } catch (e) {
      debugPrint('Error sending message: $e');
      return {
        'message': 'Sorry, I could not process your request. Please try again.',
        'suggestions': [],
      };
    }
  }

  /// POST /ai/chat/stream - Stream chat response (SSE)
  Stream<String> streamMessage(
    String message, {
    String? context,
    List<ChatMessage>? history,
  }) async* {
    final baseUrl = _client.baseUrl;
    final token = await _client.getAuthToken();

    debugPrint('🔄 Starting chat stream request...');
    debugPrint('   URL: $baseUrl/ai/chat/stream');
    debugPrint('   Message: $message');

    final request = http.Request(
      'POST',
      Uri.parse('$baseUrl/ai/chat/stream'),
    );

    request.headers.addAll({
      'Content-Type': 'application/json',
      'Accept': 'text/event-stream',
      if (token != null) 'Authorization': 'Bearer $token',
    });

    request.body = jsonEncode({
      'message': message,
      if (context != null) 'context': context,
      if (history != null) 'history': history.map((m) => m.toJson()).toList(),
    });

    try {
      final client = http.Client();
      final streamedResponse = await client.send(request);

      debugPrint('📡 Stream response status: ${streamedResponse.statusCode}');

      if (streamedResponse.statusCode != 200) {
        debugPrint('❌ Non-200 status, falling back to non-streaming');
        final body = await streamedResponse.stream.bytesToString();
        debugPrint('   Response: $body');

        // Fallback to non-streaming
        final fallbackResponse = await _client.post('/ai/chat', data: {
          'message': message,
          if (context != null) 'context': context,
          if (history != null)
            'history': history.map((m) => m.toJson()).toList(),
        });
        yield fallbackResponse.data['message'] ??
            'I could not process your request.';
        return;
      }

      await for (final chunk
          in streamedResponse.stream.transform(utf8.decoder)) {
        debugPrint('📨 Received chunk: ${chunk.length} bytes');
        // Parse SSE format: "data: {...}\n\n"
        for (final line in chunk.split('\n')) {
          if (line.startsWith('data: ')) {
            try {
              final jsonStr = line.substring(6);
              final data = jsonDecode(jsonStr);
              final content = data['content'] as String? ?? '';
              if (content.isNotEmpty) {
                yield content;
              }
              // Check if done
              if (data['done'] == true) {
                debugPrint('✅ Stream complete');
                client.close();
                return;
              }
            } catch (e) {
              debugPrint('⚠️ Parse error: $e');
            }
          }
        }
      }
      client.close();
    } catch (e) {
      debugPrint('❌ Stream error: $e');
      // Fallback to non-streaming
      try {
        final fallbackResponse = await _client.post('/ai/chat', data: {
          'message': message,
          if (context != null) 'context': context,
          if (history != null)
            'history': history.map((m) => m.toJson()).toList(),
        });
        yield fallbackResponse.data['message'] ??
            'I could not process your request.';
      } catch (e2) {
        debugPrint('❌ Fallback also failed: $e2');
        yield 'Sorry, I encountered an error while processing your request.';
      }
    }
  }
}

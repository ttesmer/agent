source .env
curl -s https://openrouter.ai/api/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $OPENROUTER_API_KEY" \
  -d '{
  "model": "minimax/minimax-m2.5",
  "messages": [
    {
      "role": "user",
      "content": "what'\''s the weather in munich? "
    }
  ],
  "tools": [
    {
      "type": "function",
      "function": {
        "name": "get_weather",
        "description": "Get the current weather in a given location",
        "parameters": 
          {
            "type": "object",
            "properties": 
              {
                "location": {"type": "string", "description": "City and state, e.g., '\''San Francisco, CA'\''"},
                "unit": {"type": "string", "enum": ["celsius", "fahrenheit"]}
              },
            "required": ["location", "unit"]
         }
       }
     }
   ],
}' | python3 -m json.tool

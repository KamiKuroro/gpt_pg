<template>
  <div id="app">
    <div class="chat-container">
      <div v-for="(message, index) in conversation" :key="index" :class="message.role">
        <div class="message-bubble">{{ message.content }}</div>
      </div>
    </div>
    <div class="input-container">
      <input v-model="userInput" @keyup.enter="sendMessage" placeholder="Type your message..." />
      <button @click="sendMessage">Send</button>
    </div>
  </div>
</template>

<script>
import axios from "axios";

export default {
  data() {
    return {
      userInput: "",
      conversation: [],
    };
  },
  methods: {
    async sendMessage() {
      if (this.userInput.trim() !== "") {
        this.conversation.push({
          role: "user",
          content: this.userInput,
        });

        try {
          const response = await axios.post("http://localhost:8080/gpt", {
            messages: this.conversation,
          });

          this.conversation.push({
            role: "assistant",
            content: response.data.choices[0].text,
          });
        } catch (error) {
          console.error("Error:", error);
        }

        this.userInput = "";
      }
    },
  },
};
</script>

<style>
body {
  font-family: Arial, sans-serif;
  margin: 0;
  padding: 0;
  display: flex;
  min-height: 100vh;
  flex-direction: column;
}

#app {
  display: flex;
  flex-direction: column;
  flex: 1;
}

.chat-container {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
}

.message-bubble {
  padding: 10px;
  border-radius: 5px;
  margin-bottom: 10px;
  display: inline-block;
}

.user .message-bubble {
  background-color: #1e88e5;
  color: white;
  align-self: flex-end;
}

.assistant .message-bubble {
  background-color: #f1f1f1;
  color: black;
}

.input-container {
  display: flex;
  padding: 20px;
  background-color: #f7f7f7;
  border-top: 1px solid #eee;
}

input {
  flex: 1;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 5px;
  outline: none;
}

button {
  margin-left: 10px;
  padding: 10px;
  background-color: #1e88e5;
  color: white;
  border: none;
  border-radius: 5px;
  cursor: pointer;
}
</style>

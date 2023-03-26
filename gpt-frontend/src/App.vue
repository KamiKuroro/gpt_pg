<template>
  <div id="app">
    <input v-model="userInput" @keyup.enter="sendMessage" placeholder="Type your message..." />
    <button @click="sendMessage">Send</button>
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

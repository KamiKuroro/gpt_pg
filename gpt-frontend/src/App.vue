<template>
  <div id="app">
    <h1>GPT-3 Text Generation</h1>
    <div class="form">
      <label for="prompt">Prompt:</label>
      <input
        type="text"
        id="prompt"
        v-model="prompt"
        @keyup.enter="fetchGeneratedText"
      />
      <button @click="fetchGeneratedText">Generate</button>
    </div>
    <div v-if="response" class="response">
      <h3>Generated Text:</h3>
      <p>{{ response }}</p>
    </div>
  </div>
</template>

<script>
import axios from "axios";

export default {
  data() {
    return {
      prompt: "",
      response: null,
    };
  },
  methods: {
    async fetchGeneratedText() {
      try {
        const res = await axios.post("http://localhost:8080/gpt", {
          prompt: this.prompt,
        });
        this.response = res.data.response;
      } catch (error) {
        console.error("Error fetching generated text:", error);
        this.response = "An error occurred while fetching the generated text.";
      }
    },
  },
};
</script>

<style>
#app {
  font-family: "Avenir", Helvetica, Arial, sans-serif;
  text-align: center;
  margin-top: 60px;
}

.form {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 8px;
}

input[type="text"] {
  padding: 4px;
  border: 1px solid #ccc;
  border-radius: 4px;
}

button {
  padding: 6px 12px;
  border: none;
  background-color: #42b983;
  color: white;
  border-radius: 4px;
  cursor: pointer;
}

button:hover {
  background-color: #2c8f73;
}

.response {
  margin-top: 40px;
}
</style>


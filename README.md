# 🤖 Mr AI - Your FREE Personal AI Assistant in the Terminal

Say hello to **Mr AI**, your no-nonsense, always-on, always-helpful terminal companion. Need code snippets, explanations, or just someone to talk CLI with? Mr AI's got your back — straight from your shell. 🐚💡

---

## ⚡ Features

- 🧠 AI-powered answers to your dev questions
- 🛠️ Code generation, debugging tips, and shell commands
- 🔄 Automatic model switching when rate limits are hit
- 🧵 Maintains conversation context throughout your session
- 🖥️ Lightweight CLI experience

---

## 🚀 Installation

```bash
# Clone the repo
git clone https://github.com/mrjxtr-dev/mr-aiCLI.git
cd mr-aiCLI

# Build the Go application
go build -o mrai

# Create a .env file with your API key
echo "OPENROUTER_API_KEY=your_api_key_here" > .env

# Run it
./mrai
```

Or set up a shortcut by adding to your shell profile:

```bash
alias mrai="/path/to/mrai"
```

---

## 🗣️ Usage

Talk to Mr AI like you would to your coding buddy:

```bash
$ mrai
You: Hey, can you explain Go interfaces?

Mr-AI: Interfaces in Go define behavior through methods. They're implemented implicitly when a type has all the required methods. This enables polymorphism without inheritance. For example, `io.Reader` defines anything that can be read from with a single `Read` method.
```

---

## 🧰 Requirements

- Go 1.24+
- OpenRouter API Key 🔑
- Internet connection 🌐

---

## 🔌 Supported Models

Mr AI uses OpenRouter to access multiple AI models with automatic failover:

- Google Gemini 2.0
- Optimus Alpha
- Llama 4 Scout
- Nvidia Llama 3.1 Nemotron Ultra

---

## 🌱 Future Plans

- [ ] Local file context for more relevant answers
- [ ] Shell command execution capabilities
- [ ] Offline mode with local models

---

## 🤝 Contributing

Pull requests welcome. Just keep it clean, useful, and not evil.

---

## 🧑‍💻 Maintained with ❤️‍🔥 by

### **Jester Lumacad** aka @[MrJxtr](https://github.com/mrjxtr)

🇵🇭 Full-Stack Software Developer from the Philippines

### 🛠️ Expertise

- API Integration & Custom Software
- Process Automation & Scripting
- Full-Stack Web Development
- AI/ML Integration & Analytics

### 💻 Tech Stack

- Languages: Python, Go, JavaScript, TypeScript

### 📫 Connect

- GitHub: [@mrjxtr](https://github.com/mrjxtr)
- LinkedIn: [Jester Lumacad](https://linkedin.com/in/mrjxtr)
- Twitter: [@mrjxtr](https://twitter.com/mrjxtr)

---

## 📜 License

MIT — Use it, fork it, improve it. Just don't make it evil.

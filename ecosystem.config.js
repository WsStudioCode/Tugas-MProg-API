module.exports = {
    apps: [
      {
        name: "go-app",
        script: "go",
        args: "run main.go",
        exec_mode: "fork",
        watch: ["main.go"],
        interpreter: "none",
      },
    ],
  };
  
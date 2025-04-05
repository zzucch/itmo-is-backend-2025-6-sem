tsParticles.load("tsparticles", {
  fpsLimit: 60,
  particles: {
    number: {
      value: 100,
      density: { enable: true, value_area: 800 },
    },
    color: { value: "#aaa" },
    shape: { type: "circle" },
    opacity: { value: 0.7, random: true },
    size: { value: 4, random: true },
    move: {
      enable: true,
      speed: 3,
      direction: "none",
      outMode: "bounce",
    },
  },
  interactivity: {
    events: {
      onHover: { enable: true, mode: "repulse" },
      onClick: { enable: true, mode: "push" },
    },
    modes: {
      repulse: { distance: 100 },
      push: { particles_nb: 4 },
    },
  },
  detectRetina: true,
});

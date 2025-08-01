import confetti from 'canvas-confetti';

export function smallButtonBurst(button: HTMLButtonElement) {
  const rect = button.getBoundingClientRect();
  const x = (rect.left + rect.width / 2) / window.innerWidth;
  const y = (rect.top + rect.height / 2) / window.innerHeight;

  confetti({
    particleCount: 40,
    startVelocity: 15,
    spread: 30,
    origin: { x, y },
    scalar: 0.6,
  });
}

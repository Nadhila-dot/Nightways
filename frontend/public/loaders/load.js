(function() {
  'use strict';

  // Load anime.js if not already loaded
  if (typeof anime === 'undefined') {
    const script = document.createElement('script');
    script.src = 'https://cdnjs.cloudflare.com/ajax/libs/animejs/3.2.1/anime.min.js';
    script.onload = initLoader;
    document.head.appendChild(script);
  } else {
    initLoader();
  }

  function initLoader() {
    // Create loader container
    const loader = document.createElement('div');
    loader.id = 'nadhi-loader';
    loader.innerHTML = `
      <div class="loader-content">
        <div class="geometric-grid">
          <div class="grid-block block-1"></div>
          <div class="grid-block block-2"></div>
          <div class="grid-block block-3"></div>
          <div class="grid-block block-4"></div>
          <div class="grid-block block-5"></div>
        </div>
        
        <div class="loader-center">
          <div class="brand-wrapper">
            <div class="brand-text">NADHI.DEV</div>
            <div class="brand-shadow">NADHI.DEV</div>
          </div>
          
          <div class="divider-container">
            <div class="divider-line"></div>
          </div>
          
          <button class="enter-btn" style="opacity: 0; pointer-events: none;">
            <span class="btn-bg"></span>
            <span class="btn-text">ENTER SITE</span>
            <span class="arrow">→</span>
          </button>
        </div>
        
        <div class="progress-bottom">
          <span class="progress-count">00</span>
        </div>
        
        <div class="frame-corners">
          <div class="corner-marker top-left"></div>
          <div class="corner-marker top-right"></div>
          <div class="corner-marker bottom-left"></div>
          <div class="corner-marker bottom-right"></div>
        </div>
        
        <div class="scanline"></div>
      </div>
    `;

    // Inject styles
    const style = document.createElement('style');
    style.textContent = `
      #nadhi-loader {
        position: fixed;
        top: 0;
        left: 0;
        width: 100vw;
        height: 100vh;
        background: #000;
        z-index: 999999;
        display: flex;
        align-items: center;
        justify-content: center;
        font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
        overflow: hidden;
      }

      #nadhi-loader * {
        box-sizing: border-box;
      }

      .loader-content {
        position: relative;
        width: 100%;
        height: 100%;
        display: flex;
        align-items: center;
        justify-content: center;
      }

      .geometric-grid {
        position: absolute;
        width: 100%;
        height: 100%;
        top: 0;
        left: 0;
        overflow: hidden;
      }

      .grid-block {
        position: absolute;
        background: #fff;
        opacity: 0;
      }

      .block-1 {
        width: 30vw;
        height: 100vh;
        left: -30vw;
        top: 0;
      }

      .block-2 {
        width: 100vw;
        height: 25vh;
        top: -25vh;
        left: 0;
      }

      .block-3 {
        width: 20vw;
        height: 100vh;
        right: -20vw;
        top: 0;
      }

      .block-4 {
        width: 100vw;
        height: 15vh;
        bottom: -15vh;
        left: 0;
      }

      .block-5 {
        width: 35vw;
        height: 35vh;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%) scale(0);
      }

      .loader-center {
        position: relative;
        z-index: 10;
        text-align: center;
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 0; /* MAIN GAP - controls spacing between all center elements */
      }

      .progress-bottom {
        position: absolute;
        bottom: 3rem;
        right: 3rem;
        font-family: 'Courier New', monospace;
        color: #fff;
        z-index: 20;
      }

      .progress-count {
        font-size: 3rem;
        font-weight: 900;
        letter-spacing: 0.2em;
        line-height: 1;
        font-feature-settings: 'tnum';
        text-shadow: 0 0 30px rgba(255, 255, 255, 0.4);
        opacity: 0.7;
      }

      .brand-wrapper {
        position: relative;
        display: inline-block;
      }

      .brand-text {
        font-size: 5.5rem;
        font-weight: 900;
        letter-spacing: -0.03em;
        color: #fff;
        text-transform: uppercase;
        position: relative;
        z-index: 2;
        opacity: 0;
        transform: translateY(30px);
        margin-bottom: 0rem; /* MARGIN AFTER BRAND TEXT - space before divider */
      }

      .brand-shadow {
        font-size: 5.5rem;
        font-weight: 900;
        letter-spacing: -0.03em;
        color: transparent;
        text-transform: uppercase;
        position: absolute;
        top: 8px;
        left: 8px;
        z-index: 1;
        -webkit-text-stroke: 1px rgba(255, 255, 255, 0.3);
        opacity: 0;
      }

      .divider-container {
        position: relative;
        width: 500px;
        height: 4px;
        display: flex;
        align-items: center;
        justify-content: center;
        margin: 0 0 1.5rem 0; /* MARGIN AFTER DIVIDER - space before button */
      }

      .divider-line {
        width: 0;
        height: 4px;
        background: #fff;
        position: absolute;
        box-shadow: 0 0 20px rgba(255, 255, 255, 0.6);
      }

      .enter-btn {
        background: transparent;
        color: #fff;
        border: 3px solid #fff;
        padding: 1.75rem 5rem;
        font-size: 1.1rem;
        font-weight: 900;
        letter-spacing: 0.4em;
        cursor: pointer;
        display: flex;
        align-items: center;
        gap: 2rem;
        text-transform: uppercase;
        position: relative;
        overflow: hidden;
        transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
        box-shadow: 0 0 30px rgba(255, 255, 255, 0.2);
      }

      .btn-bg {
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        background: #fff;
        transform: translateX(-100%);
        transition: transform 0.5s cubic-bezier(0.4, 0, 0.2, 1);
        z-index: 0;
      }

      .enter-btn:hover .btn-bg,
      .enter-btn:active .btn-bg {
        transform: translateX(0);
      }

      .enter-btn:hover,
      .enter-btn:active {
        color: #000;
        box-shadow: 0 0 50px rgba(255, 255, 255, 0.6);
        transform: scale(1.05);
      }

      .enter-btn:active {
        transform: scale(0.98);
      }

      .btn-text {
        position: relative;
        z-index: 1;
      }

      .arrow {
        font-size: 1.5rem;
        transition: transform 0.3s ease;
        position: relative;
        z-index: 1;
      }

      .enter-btn:hover .arrow {
        transform: translateX(8px);
      }

      .frame-corners {
        position: absolute;
        width: 100%;
        height: 100%;
        top: 0;
        left: 0;
        pointer-events: none;
      }

      .corner-marker {
        position: absolute;
        width: 50px;
        height: 50px;
        opacity: 0;
      }

      .top-left {
        top: 2rem;
        left: 2rem;
        border-top: 3px solid #fff;
        border-left: 3px solid #fff;
      }

      .top-right {
        top: 2rem;
        right: 2rem;
        border-top: 3px solid #fff;
        border-right: 3px solid #fff;
      }

      .bottom-left {
        bottom: 2rem;
        left: 2rem;
        border-bottom: 3px solid #fff;
        border-left: 3px solid #fff;
      }

      .bottom-right {
        bottom: 2rem;
        right: 2rem;
        border-bottom: 3px solid #fff;
        border-right: 3px solid #fff;
      }

      .scanline {
        position: absolute;
        width: 100%;
        height: 2px;
        background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.8), transparent);
        top: -2px;
        opacity: 0;
        filter: blur(1px);
      }

      @media (max-width: 768px) {
        .brand-text, .brand-shadow {
          font-size: 3.5rem;
        }
        .progress-count {
          font-size: 2rem;
        }
        .progress-bottom {
          bottom: 2rem;
          right: 2rem;
        }
        .enter-btn {
          padding: 1.25rem 3rem;
          font-size: 0.9rem;
        }
        .divider-container {
          width: 300px;
        }
      }
    `;

    document.head.appendChild(style);
    document.body.appendChild(loader);

    // Start animations
    startAnimation();
  }

  function startAnimation() {
    const loader = document.getElementById('nadhi-loader');
    const progressCount = loader.querySelector('.progress-count');
    const progressBar = loader.querySelector('.progress-bar');
    const enterBtn = loader.querySelector('.enter-btn');

    // Scanline animation
    anime({
      targets: '#nadhi-loader .scanline',
      translateY: ['0vh', '100vh'],
      opacity: [0, 0.6, 0],
      duration: 2000,
      easing: 'linear',
      loop: 2
    });

    // Geometric blocks entrance - MORE AGGRESSIVE
    const timeline = anime.timeline({
      easing: 'easeInOutExpo'
    });

    timeline
      .add({
        targets: '#nadhi-loader .block-1',
        opacity: [0, 0.95, 0],
        left: ['-30vw', '-5vw', '40vw'],
        duration: 1400,
        easing: 'easeInOutQuint'
      })
      .add({
        targets: '#nadhi-loader .block-2',
        opacity: [0, 0.95, 0],
        top: ['-25vh', '-3vh', '30vh'],
        duration: 1200,
        easing: 'easeInOutQuint'
      }, '-=1000')
      .add({
        targets: '#nadhi-loader .block-3',
        opacity: [0, 0.95, 0],
        right: ['-20vw', '-3vw', '25vw'],
        duration: 1300,
        easing: 'easeInOutQuint'
      }, '-=900')
      .add({
        targets: '#nadhi-loader .block-4',
        opacity: [0, 0.95, 0],
        bottom: ['-15vh', '-2vh', '20vh'],
        duration: 1100,
        easing: 'easeInOutQuint'
      }, '-=800')
      .add({
        targets: '#nadhi-loader .block-5',
        opacity: [0, 0.9, 0],
        scale: [0, 2, 0],
        rotate: [0, 180],
        duration: 1600,
        easing: 'easeInOutQuint'
      }, '-=700');

    // Progress counter (bottom right) - LUXURY EASING
    anime({
      targets: { value: 0 },
      value: 100,
      duration: 3200,
      easing: 'easeInOutExpo',
      update: function(anim) {
        const val = Math.round(anim.animations[0].currentValue);
        progressCount.textContent = val.toString().padStart(2, '0');
      }
    });

    // Brand text with shadow reveal - STAGGERED LUXURY
    anime({
      targets: '#nadhi-loader .brand-shadow',
      opacity: [0, 0.6],
      translateY: [40, 0],
      translateX: [-10, 0],
      duration: 1400,
      delay: 1200,
      easing: 'easeOutExpo'
    });

    anime({
      targets: '#nadhi-loader .brand-text',
      opacity: [0, 1],
      translateY: [50, 0],
      duration: 1600,
      delay: 1300,
      easing: 'easeOutExpo'
    });

    // Divider line expand - SMOOTH LUXURY GROWTH
    anime({
      targets: '#nadhi-loader .divider-line',
      width: ['0px', '500px'],
      duration: 1800,
      delay: 2400,
      easing: 'easeInOutExpo'
    });

    // Corner markers with stagger - PREMIUM REVEAL
    anime({
      targets: '#nadhi-loader .corner-marker',
      opacity: [0, 1],
      scale: [0.3, 1],
      duration: 1000,
      delay: anime.stagger(150, { start: 2600 }),
      easing: 'easeOutExpo'
    });

    // Enter button reveal - LUXURY ENTRANCE
    setTimeout(() => {
      enterBtn.style.opacity = '1';
      enterBtn.style.pointerEvents = 'auto';
      
      anime({
        targets: enterBtn,
        opacity: [0, 1],
        translateY: [40, 0],
        scale: [0.92, 1],
        duration: 1200,
        easing: 'easeOutExpo'
      });
    }, 3400);

    // Enter button click handler (or any key press) - LUXURY EXIT
    const exitLoader = () => {
      // Disable further interactions
      enterBtn.style.pointerEvents = 'none';
      document.removeEventListener('keydown', handleKeyPress);

      // Premium exit sequence
      anime({
        targets: '#nadhi-loader .loader-center',
        opacity: [1, 0],
        scale: [1, 0.95],
        duration: 600,
        easing: 'easeInExpo'
      });

      anime({
        targets: '#nadhi-loader .corner-marker',
        opacity: [1, 0],
        scale: [1, 0.8],
        duration: 500,
        delay: anime.stagger(50),
        easing: 'easeInExpo'
      });

      anime({
        targets: '#nadhi-loader .progress-bottom',
        opacity: [1, 0],
        translateX: [0, 30],
        duration: 500,
        easing: 'easeInExpo'
      });

      // Final fade with blocks flash
      /*
      anime({
        targets: '#nadhi-loader .grid-block',
        opacity: [0, 0.3, 0],
        duration: 800,
        delay: 300,
        easing: 'easeInOutQuad'
      });*/

      anime({
        targets: '#nadhi-loader',
        opacity: [1, 0],
        duration: 800,
        delay: 500,
        easing: 'easeInExpo',
        complete: () => {
          loader.remove();
        }
      });
    };

    const handleKeyPress = (e) => {
      if (enterBtn.style.opacity === '1') {
        exitLoader();
      }
    };

    enterBtn.addEventListener('click', exitLoader);
    document.addEventListener('keydown', handleKeyPress);
  }
})();
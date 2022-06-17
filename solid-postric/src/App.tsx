import type { Component } from 'solid-js';

import classnames from "classnames";

import logo from './assets/logo.png';
import banner from './assets/login/banner.png';
import styles from './App.module.css';

const App: Component = () => {
  return (
    <div class={styles.container}>
      <header class={styles.header}>
        <img src={logo} alt="Logo" class={styles.logo} />

        <section class={styles.navigation}>
          <span><a href="#">Home</a></span>
          <span><a href="#">Testimony</a></span>
          <span><a href="#">About Us</a></span>
          <span><a href="#">Contact Us</a></span>
        </section>

        <button class={styles.button}>Sign Up</button>
      </header>
      <main class={styles.main}>
        <section class={styles.bannerWrapper}>
          <img src={banner} alt="Banner" class={styles.banner} />
        </section>

        <section class={styles.signup}>
          <span>Sign Up</span>
          <p>If Already Signed Up then Login In</p>

          <input type="text" name="name" class={styles.input} placeholder="Name" />
          <input type="text" name="email" class={styles.input} placeholder="Email" />
          <input type="text" name="password" class={styles.input} placeholder="Password" />
          <input type="text" name="phone" class={styles.input} placeholder="Phone" />

          <button class={classnames(styles.button, styles.fullWidth)}>Sign Up</button>
          <p>Or Sign Up Using</p>
        </section>
      </main>
    </div>
  );
};

export default App;

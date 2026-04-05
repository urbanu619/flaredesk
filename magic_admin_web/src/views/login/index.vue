<template>
  <div class="login-container flx-center">
    <div class="login-box">
      <SwitchDark class="dark"/>
      <div class="login-left">
        <img class="login-left-img" src="@/assets/images/login_left.png" alt="login"/>
      </div>
      <div class="login-form">
        <div class="login-logo">
          <h2 class="logo-text">{{ title }}</h2>
        </div>
        <LoginForm/>
      </div>
    </div>
  </div>
</template>

<script setup  name="login">
import LoginForm from "./components/LoginForm.vue";
import SwitchDark from "@/components/SwitchDark/index.vue";

const title = import.meta.env.VITE_GLOB_APP_TITLE;
</script>

<style scoped lang="scss">
@use "./index.scss";
</style>

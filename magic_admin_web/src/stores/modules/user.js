import {defineStore} from "pinia";


export const useUserStore = defineStore({
  id: "userInfo",
  state: () => ({
    token: "",
    userInfo: {
      avatar: "",
      isSetGoogleAuth: false,
      nickname: ""
    }
  }),

  actions: {
    // Set Token
    setToken(token) {
      this.token = token;
    },
    // Set setUserInfo
    setUserInfo(userInfo) {
      this.userInfo = userInfo;
    }
  },
  persist: true
});

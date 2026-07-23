import i18n from 'i18next'
import { initReactI18next } from 'react-i18next'

const resources = {
  en: {
    translation: {
      appName: 'YCTF',
      challenges: 'Challenges',
      scoreboard: 'Scoreboard',
      profile: 'Profile',
      admin: 'Admin',
      login: 'Login',
      register: 'Register',
      logout: 'Logout',
      username: 'Username',
      email: 'Email',
      password: 'Password',
      submit: 'Submit',
      flag: 'Flag',
      points: 'Points',
      solves: 'Solves',
      rank: 'Rank',
      team: 'Team',
      score: 'Score',
    },
  },
  zh: {
    translation: {
      appName: 'YCTF',
      challenges: '题目',
      scoreboard: '排行榜',
      profile: '个人中心',
      admin: '管理',
      login: '登录',
      register: '注册',
      logout: '退出',
      username: '用户名',
      email: '邮箱',
      password: '密码',
      submit: '提交',
      flag: 'Flag',
      points: '分值',
      solves: '解题数',
      rank: '排名',
      team: '队伍',
      score: '分数',
    },
  },
}

i18n.use(initReactI18next).init({
  resources,
  lng: 'zh',
  fallbackLng: 'en',
  interpolation: { escapeValue: false },
})

export default i18n

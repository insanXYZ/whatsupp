import axios from "axios";
import 'dotenv/config';

export const API = axios.create({
  baseURL: process.env.BASE_URL_API,
  withCredentials: true,
});

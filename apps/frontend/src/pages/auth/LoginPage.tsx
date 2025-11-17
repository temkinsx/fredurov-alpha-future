import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import useSignIn from "react-auth-kit/hooks/useSignIn";
import useFetch from "../../hooks/useFetch";

import { Sparkles, Shield, Zap } from "lucide-react";
import type { UserInfo } from "../../lib/types";

const LoginPage = () => {
  const API_BASE = "http://localhost:8080";

  const signIn = useSignIn();
  const navigate = useNavigate();

  const [_, setUserInfo] = useState<UserInfo | null>(null);
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");


  const { fetching: fetchUserInfo, isPending: isUserInfoPending } = useFetch(async () => {
    const response = await fetch(`${API_BASE}/login`, {
      method: "POST",
      headers: { 
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ email, password }),
    });

    if (!response.ok) throw new Error("Invalid credentials");

    const data: UserInfo = await response.json();
    setUserInfo(data);

    const success = signIn({
      auth: {
        token: `${email}:${password}`,
        type: "Basic",
      },
      userState: { email: data.email },
    });

    if (success) navigate("/");
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    fetchUserInfo();
  };

  return (
    <div className="min-h-screen flex items-center justify-center px-6 py-12 bg-[#ffffff]">
      <div className="w-full max-w-6xl grid grid-cols-2 gap-12 items-end">

        <section className="space-y-8">
          <div className="space-y-3">
            <div className="inline-block bg-[#ef3124] px-4 py-2 rounded-md">
              <span className="text-white text-sm">Альфа-Банк</span>
            </div>

            <h1 className="text-4xl font-semibold">Alfa Copilot</h1>

            <p className="text-blue-700">AI-помощник для вашего бизнеса</p>

            <p className="text-gray-600">
              Войдите в систему для доступа к персональному AI-ассистенту
            </p>
          </div>

          <div className="space-y-4 gap-4">
            {[
              { icon: Sparkles, title: "AI-Аналитика", desc: "Умный анализ финансовых данных" },
              { icon: Shield, title: "Безопасность", desc: "Банковская защита данных" },
              { icon: Zap, title: "Мгновенно", desc: "Быстрые ответы 24/7" },
            ].map((item, i) => (
              <div
                key={i}
                className={`
                  border-t border-t-gray-100 rounded-xl p-4 bg-white h-auto
                  flex flex-row gap-5 items-center shadow-md
                  ${
                    i === 1
                      ? "shadow-[#9933ff3c]"
                      : i === 2
                      ? "shadow-[#ef31243c]" 
                      : "shadow-[#a9ff003c]"
                  }
                `}
              >
                <div className="flex items-center justify-center w-12 h-12">
                  <item.icon className="w-6 h-6 text-red-500" />
                </div>

                <div>
                  <h3 className="font-medium">{item.title}</h3>
                  <p className="text-sm text-gray-500">{item.desc}</p>
                </div>
              </div>
            ))}
          </div>
        </section>

        <section>
          <form
            onSubmit={handleSubmit}
            className="bg-white border-t border-t-gray-100 rounded-2xl p-8 space-y-6 h-auto shadow-lg shadow-[#9933ff3c]"
          >
            <div>
              <h2 className="text-2xl font-semibold">Вход в систему</h2>
              <p className="text-sm text-gray-500">Используйте корпоративные данные</p>
            </div>

            <div className="space-y-2">
              <label className="text-sm text-gray-600">Email</label>
              <input
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="your.email@alfabank.ru"
                className="w-full px-4 py-3 border rounded-xl outline-none"
                required
              />
            </div>

            <div className="space-y-2">
              <label className="text-sm text-gray-600">Пароль</label>
              <input
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="••••••••"
                className="w-full px-4 py-3 border rounded-xl outline-none"
                required
              />
            </div>

            <div className="flex items-center justify-between text-sm">
              <label className="flex items-center gap-2">
                <input type="checkbox" className="w-4 h-4" />
                Запомнить меня
              </label>

              <button type="button" className="text-blue-600">
                Забыли пароль?
              </button>
            </div>

            <button
              type="submit"
              className="w-full py-3 rounded-xl bg-[#ef3124] text-white"
            >
              {isUserInfoPending ? "Загрузка..." : "Войти в Alfa Copilot"}
            </button>

            <p className="text-center text-sm text-gray-500">
              Нет доступа? Обратитесь к{" "}
              <button className="text-blue-600">администратору</button>
            </p>
          </form>
        </section>

      </div>
    </div>
  );
};

export default LoginPage;
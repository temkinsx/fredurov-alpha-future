import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import useSignIn from "react-auth-kit/hooks/useSignIn";
import useFetch from "../../hooks/useFetch";

import type { UserInfo } from "../../lib/types";

const LoginPage = () => {

  const API_BASE = "https://maxbot-withoutdocker.onrender.com";

  const signIn = useSignIn();
  const navigate = useNavigate();

  const [userInfo, setUserInfo] = useState<UserInfo | null>(null);
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  const { fetching: fetchUserInfo, isPending: isUserInfoPending } = useFetch(async () => {
    const response = await fetch(`${API_BASE}/login/`);

    if (!response.ok) throw new Error("Invalid credentials");

    const data: UserInfo = await response.json();
    setUserInfo(data);

    const success = signIn({
      auth: {
        token: `${email}:${password}`,
        type: "Basic",
      },
      userState: {
        email: data.email,
      },
    });

    if (success) {
      navigate("/");
    }
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault(); // prevents page reload
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

            <p className="text-blue-700">
              AI-помощник для вашего бизнеса
            </p>

            <p className="text-gray-600">
              Войдите в систему для доступа к персональному AI-ассистенту
            </p>
          </div>

          <div className="space-y-4">
            {[
              {
                title: "AI-Аналитика",
                desc: "Умный анализ финансовых данных",
              },
              {
                title: "Безопасность",
                desc: "Банковская защита данных",
              },
              {
                title: "Мгновенно",
                desc: "Быстрые ответы 24/7",
              },
            ].map((item, i) => (
              <div
                key={i}
                className="border rounded-xl p-4 bg-white h-auto"
              >
                <h3 className="font-medium">{item.title}</h3>
                <p className="text-sm text-gray-500">{item.desc}</p>
              </div>
            ))}
          </div>
        </section>

        <section>
          <form className="bg-white border rounded-2xl p-8 space-y-6 h-auto">
            <div>
              <h2 className="text-2xl font-semibold">Вход в систему</h2>
              <p className="text-sm text-gray-500">
                Используйте корпоративные данные
              </p>
            </div>

            <div className="space-y-2">
              <label className="text-sm text-gray-600">Email</label>
              <input
                type="email"
                placeholder="your.email@alfabank.ru"
                className="w-full px-4 py-3 border rounded-xl outline-none"
              />
            </div>

            <div className="space-y-2">
              <label className="text-sm text-gray-600">Пароль</label>
              <input
                type="password"
                placeholder="••••••••"
                className="w-full px-4 py-3 border rounded-xl outline-none"
              />
            </div>

            <div className="flex items-center justify-between text-sm">
              <label className="flex items-center gap-2">
                <input type="checkbox" className="w-4 h-4" />
                Запомнить меня
              </label>
              <button className="text-blue-600">Забыли пароль?</button>
            </div>

            <Link to="/">
              <button
                type="submit"
                className="w-full py-3 rounded-xl bg-[#ef3124] text-white"
              >
                Войти в Alfa Copilot
              </button>
            </Link>

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
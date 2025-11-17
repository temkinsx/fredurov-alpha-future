import type { Message } from "../../../../types.ts";
import MessageNode from "./MesageNode.tsx";


const mockMessages: Message[] = [
    {
        content: "Добрый день! Я Alfa Copilot — ваш AI-помощник для бизнеса. Чем могу помочь?",
        isAnswer: true
    },
    {
        content: "Покажи финансовую аналитику за последний квартал",
        isAnswer: false
    },
    {
        content: "Анализирую данные за Q3 2024... Ваша выручка выросла на 23% по сравнению с прелылушим кварталом. Основной рост в сегменте В2В",
        isAnswer: true
    }
]

function MainChat() {

    return (
        <div className="flex flex-col gap-4 mt-10">
            {mockMessages.map(msg => (
                <MessageNode content={ msg.content } isAnswer={ msg.isAnswer } />
            ))}
        </div>
    );
}

export default MainChat;
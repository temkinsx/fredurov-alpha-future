

function ChatHeader() {


    return (
        <header className="flex gap-4 p-5 border-b border-b-gray-100" style={{gridColumn: "2/3", gridRow: "1/2"}}>
            <a className="bg-red-500 p-3.5 rounded-2xl
            hover:bg-red-600 transition-all"
               style={{color: "white"}}
               href="https://alfabank.ru"
               target="_blank"
            >
                Альфа-Банк
            </a>
            <h1 className="flex flex-col">
                <span className="text-2xl">Alpha Copilot</span>
                <span className="text-gray-500 text-sm">AI-помощник для бизнеса</span>
            </h1>
        </header>
    );
}

export default ChatHeader;
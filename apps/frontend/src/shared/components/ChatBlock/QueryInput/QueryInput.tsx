import { useState } from "react";
import { Send } from "lucide-react";


function QueryInput() {
    const [queryValue, setQueryValue] = useState<string>("");

    const handleSendQuery = () => {

    }

    return (
        <div
            className="flex justify-center items-center p-3
            bg-white border-t border-t-gray-100"
            style={{gridColumn: "2/3", gridRow: "10/11"}}
        >
            <div className="flex flex-col p-3
            bg-white border border-gray-100 shadow-lg shadow-gray-200 rounded-2xl">
                <textarea
                    className="resize-none border-none outline-none"
                    value={ queryValue }
                    onChange={ e => setQueryValue(e.target.value) }
                    placeholder="Задайте вопрос Alpha Copilot..."
                    cols={100}
                    rows={3}
                >
                </textarea>
                <div className="bg-gray-100 p-2 rounded-2xl
                self-end mt-2 hover:bg-gray-200 transition-all">
                    <button
                        className="flex justify-center items-center gap-2 text-gray-500"
                        onClick={ handleSendQuery }
                    >
                        <Send className="h-5 w-5" />
                        Отправить
                    </button>
                </div>
            </div>
        </div>
    );
}

export default QueryInput;
import type { Message } from "../../../../types.ts";


function MessageNode( {content, isAnswer}: Message ) {

    return (
        <div
            className={` ${isAnswer ? 
                "bg-white shadow-lg shadow-gray-200 self-start" 
                : "bg-red-500 text-white shadow-lg shadow-red-200 self-end"}
                p-4 max-w-3/4 rounded-2xl`}
        >
            <p>{ content }</p>
        </div>
    );
}

export default MessageNode;
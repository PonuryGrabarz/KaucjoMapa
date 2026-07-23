import Styles from "./alert-item.module.sass";
import IconifyIcon from "@components/iconify_icon/iconify-icon";

import type { AlertItem as AlertItemData } from "../../lib/alerts-mock-data";

type Props = AlertItemData;

export default function AlertItem({
    id,
    icon,
    title,
    description,
    timestamp,
    dateTime,
    unread,
}: Props) {
    return (
        <li class={Styles["alert-item"]}>
            <article
                class={Styles["alert-item__card"]}
                data-unread={unread ? "true" : "false"}
                aria-labelledby={`${id}-title`}
            >
                <span class={Styles["alert-item__icon"]} aria-hidden="true">
                    <IconifyIcon icon={icon} />
                </span>

                <div class={Styles["alert-item__content"]}>
                    <h2 id={`${id}-title`}>{title}</h2>
                    <p>{description}</p>
                    <footer>
                        <time dateTime={dateTime}>{timestamp}</time>
                    </footer>
                </div>

                {unread && (
                    <span
                        class={Styles["alert-item__unread"]}
                        role="img"
                        aria-label="Nieprzeczytane powiadomienie"
                    />
                )}
            </article>
        </li>
    );
}

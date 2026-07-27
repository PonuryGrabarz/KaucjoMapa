import Styles from "./alerts-list.module.sass";
import AlertItem from "../alert_item/alert-item";

import type { AlertItem as AlertItemData } from "../../lib/alerts-mock-data";

interface Props {
    id: string;
    items: AlertItemData[];
    allRead: boolean;
}

export default function AlertsList({ id, items, allRead }: Props) {
    return (
        <ul id={id} class={Styles["alerts-list"]} aria-label="Powiadomienia o punktach zwrotu">
            {items.map((item) => (
                <AlertItem key={item.id} {...item} unread={item.unread && !allRead} />
            ))}
        </ul>
    );
}

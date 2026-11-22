import type { ReactNode } from 'react'
import AppText from '../../AppText/ui/AppText.tsx'
import AppTitle from '../../AppTitle/ui/AppTitle.tsx'
import styles from './PageHeader.module.css'

export interface PageHeaderProps {
  /** Хлебные крошки, например «Главная / Обзор». */
  breadcrumbs: string
  /** Заголовок страницы. */
  title: string
  /** Подзаголовок/описание страницы. */
  description: string
  /** Опциональный слот справа от заголовка (действия страницы). */
  actions?: ReactNode
}

export function PageHeader({ breadcrumbs, title, description, actions }: PageHeaderProps) {
  return (
    <div className={styles.root}>
      <AppText className={styles.crumbs}>{breadcrumbs}</AppText>
      <div className={styles.titleRow}>
        <AppTitle className={styles.heading}>{title}</AppTitle>
        {actions != null && <div className={styles.actions}>{actions}</div>}
      </div>
      <AppText className={styles.subtitle}>{description}</AppText>
    </div>
  )
}

export default PageHeader

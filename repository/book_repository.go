package repository

import (
	"database/sql"
	"fmt"
	"go-api/model"

	"github.com/gin-gonic/gin"
)

type BookRepository struct {
	connection *sql.DB
}

func NewBookRepository (connection *sql.DB) BookRepository {
	return BookRepository{
		connection: connection,
	}
}

func (pr *BookRepository) GetAllBooks() ([]model.Book, error){
	query := "SELECT * FROM book"
	rows, err := pr.connection.Query(query)
	if(err != nil){
		fmt.Println(err)
		return []model.Book{}, err
	}

	var AllbookList []model.Book
	var bookObj model.Book

	for rows.Next(){
		err = rows.Scan(
			&bookObj.Id,
			&bookObj.Title,
			&bookObj.Synopsis,
			&bookObj.Price,
			&bookObj.Amount,
			&bookObj.Author_id)

		if(err != nil){
			fmt.Println(err)
			return []model.Book{}, err
		}

		AllbookList = append(AllbookList, bookObj)
	}

	rows.Close()

	return AllbookList, nil
}

func (pr *BookRepository) GetBooks(c *gin.Context) ([]model.Book, error){
	title := c.Query("title")
	genre := c.Query("genre")
	author := c.Query("author")

	query := "SELECT * FROM book WHERE 1=1"
	var args []interface{}
	argIdx := 1

	//filtro por title
	if title != ""{
		query += fmt.Sprintf(" AND title ILIKE $%d", argIdx)
		args = append(args, "%"+title+"%")
		argIdx++
	}

	//filtro do genero
	if genre != ""{
		query += fmt.Sprintf("  AND id IN (SELECT fk_book_id FROM genre_book gb JOIN genres g ON gb.fk_genre_id = g.id WHERE g.name ILIKE $%d)", argIdx)
		args = append(args, "%"+genre+"%")
		argIdx++
	}

	//filtro do autor
	if author != ""{
		query += fmt.Sprintf(" AND author_id IN (SELECT id FROM authors WHERE name ILIKE $%d)", argIdx)
		args = append(args, "%"+author+"%")
		argIdx++
	}

	fmt.Println("Query gerada: ", query)

	stmt, err := pr.connection.Prepare(query)
	if err != nil {
		fmt.Println("Erro na preparação da query: ", err)
		return []model.Book{}, err
	}

	defer stmt.Close()

	rows, err := stmt.Query(args...)
	if err != nil {
		fmt.Println("Erro na query: ", err)
		return []model.Book{}, err
	}

	defer rows.Close()

	var bookList []model.Book

	for rows.Next() {
		var bookObj model.Book

		err = rows.Scan(
			&bookObj.Id,
			&bookObj.Title,
			&bookObj.Synopsis,
			&bookObj.Price,
			&bookObj.Amount,
			&bookObj.Author_id,
		)

		if err != nil {
			fmt.Println("Erro no scan das linas: ", err)
			return []model.Book{}, err
		}

		bookList = append(bookList, bookObj)
	}

	return bookList, nil
}

func (pr *BookRepository) DeleteBook(c *gin.Context) (string, error){
	id := c.Param("id")
	title := c.Param("title")

	query, err := pr.connection.Prepare("DELETE FROM book WHERE id = $1")

	if err != nil {
		fmt.Println(err)
		return "",err
	}

	defer query.Close()

	result, err := query.Exec(id)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	// Verifica se alguma linha foi afetada
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	// Se nenhuma linha foi afetada, significa que o livro com o ID fornecido não existe
	if rowsAffected == 0 {
		return "",fmt.Errorf("nenhum livro encontrado com o id %s", id)
	}

	return fmt.Sprintf("O livro %s foi deletado com sucesso!", title), nil

}

func (pr *BookRepository) CreateBook(book model.Book) (string, error){
	var authorID string
	err := pr.connection.QueryRow("SELECT id FROM author WHERE name = $1", book.Author_id).Scan(&authorID)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("autor não encontrado com o nome: %s", book.Author_id)
		}
		return "", err
	}

	query, err := pr.connection.Prepare("INSERT INTO book" +
							"(title, synopsis, price, amount, author_id)" +
							"VALUES ($1, $2, $3 , $4, $5) RETURNING id")

	if err != nil {
			fmt.Println(err)
			return "",err
	}
	defer query.Close()
	var id string

	err = query.QueryRow(book.Title, book.Synopsis, book.Price, book.Amount, authorID).Scan(&id)
		if err != nil {
			fmt.Println(err)
			return "", err
	}
	
	return id, nil
}

func (pr *BookRepository) UpdateBook(book model.Book) error {
	// Prepara a consulta SQL para atualização
	query, err := pr.connection.Prepare("UPDATE book SET title = $1, synopsis = $2, price = $3, amount = $4, author_id = $5 WHERE id = $6")

	if err != nil {
			fmt.Println("Erro ao preparar a consulta:", err)
			return err
	}
	// Certifique-se de fechar a consulta após o uso
	defer query.Close()

	// Executa a consulta com os parâmetros do livro
	_, err = query.Exec(book.Title, book.Synopsis, book.Price, book.Amount, book.Author_id, book.Id)
	if err != nil {
			fmt.Println("Erro ao executar a consulta:", err)
			return err
	}

	return nil 
}
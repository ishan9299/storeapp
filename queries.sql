create table manufacturer (
    id int not null auto_increment,
    name varchar(255) not null,
    address varchar(255) not null,
    primary key(id)
);

insert into manufacturer (name, address) values("Ashok Leyland", "Chennai")
insert into manufacturer (name, address) values("Hero Honda Motors", "New Delhi");
insert into manufacturer (name, address) values("Godrej Group", "Mumbai");

create table products (
    id int not null auto_increment,
    name varchar(255) not null,
    quantity int not null,
    manufacid int,
    primary key (id),
    foreign key (manufacid) references manufacturer(id)
);

insert into products (name, quantity, manufacid) values("Aashirvaad Atta", "5", "4");
insert into products (name, quantity, manufacid) values("Bingo! Cheese Nachos", "7", "4");
insert into products (name, quantity, manufacid) values("Cinthol Soap", "3", "3");
insert into products (name, quantity, manufacid) values("Honey", "9", "8");

select p.id, p.name, p.quantity, m.name from products p inner join manufacturer m on p.manufacid = m.id;
